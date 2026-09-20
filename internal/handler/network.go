package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// IPInfo IP 归属信息，字段名与前端 IpQuery.vue 约定一致
type IPInfo struct {
	Status      string  `json:"status"`
	Message     string  `json:"message,omitempty"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

const (
	ipUserAgent    = "mytools/1.0"
	ipHTTPTimeout  = 8 * time.Second
	ipQueryTimeout = 6 * time.Second
	ipMaxBodyBytes = 1 << 20
)

type NetworkHandler struct {
	client *http.Client
}

func NewNetworkHandler() *NetworkHandler {
	return &NetworkHandler{
		client: &http.Client{
			Timeout: ipHTTPTimeout,
		},
	}
}

// ipSource 一个 IP 查询数据源
type ipSource struct {
	name  string
	build func(query string) string
	parse func([]byte) (*IPInfo, error)
}

// ipSources 多个数据源并发查询，最先成功返回的胜出。
// ip.sb / ipwho.is 在海外、字段最全；vore.top 在国内、返回中文，用来兜底。
var ipSources = []ipSource{
	{name: "ip.sb", build: buildIPSBURL, parse: parseIPSB},
	{name: "vore.top", build: buildVoreURL, parse: parseVore},
	{name: "ipwho.is", build: buildIPWhoURL, parse: parseIPWho},
}

// GetIpInfo 查询 IP 归属信息，ip 为空时查询本机公网 IP。
func (a *NetworkHandler) GetIpInfo(ip string) (*IPInfo, error) {
	query, failMsg := resolveQuery(ip)
	if failMsg != "" {
		return &IPInfo{Status: "error", Message: failMsg}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), ipQueryTimeout)
	defer cancel()

	type outcome struct {
		info   *IPInfo
		source string
		err    error
	}
	results := make(chan outcome, len(ipSources))

	for _, src := range ipSources {
		go func(src ipSource) {
			info, err := fetchFromSource(ctx, a.client, src, query)
			results <- outcome{info: info, source: src.name, err: err}
		}(src)
	}

	timeouts := 0
	for range ipSources {
		select {
		case out := <-results:
			if out.err == nil && out.info != nil {
				return out.info, nil
			}
			log.Printf("[ipinfo] 数据源 %s 失败: %v", out.source, out.err)
			var netErr net.Error
			if errors.As(out.err, &netErr) && netErr.Timeout() {
				timeouts++
			}
		case <-ctx.Done():
			log.Printf("[ipinfo] 整体查询超时（%s）", ipQueryTimeout)
			return &IPInfo{Status: "error", Message: "查询超时，请检查网络后重试"}, nil
		}
	}

	if timeouts == len(ipSources) {
		return &IPInfo{Status: "error", Message: "查询超时，请检查网络后重试"}, nil
	}
	return &IPInfo{Status: "error", Message: "所有查询服务均不可用，请稍后重试"}, nil
}

// fetchFromSource 向单个数据源发起一次查询并解析结果
func fetchFromSource(ctx context.Context, client *http.Client, src ipSource, query string) (*IPInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, src.build(query), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", ipUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, ipMaxBodyBytes))
	if err != nil {
		return nil, err
	}
	return src.parse(body)
}

// resolveQuery 校验并规范化查询条件。
// 空输入表示查询本机（返回空串，由调用方拼成不带 IP 的地址）；
// 域名会先解析成 IP —— 各数据源基本都只接受 IP 字面量。
// 第二个返回值非空表示校验失败，内容为给用户的提示。
func resolveQuery(input string) (string, string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", ""
	}
	if net.ParseIP(input) != nil {
		return input, ""
	}

	addrs, err := net.LookupIP(input)
	if err != nil || len(addrs) == 0 {
		return "", "无法解析该地址，请输入有效的 IP 地址或域名"
	}
	// 优先取 IPv4：多数数据源对 IPv6 给不出省市信息
	picked := addrs[0]
	for _, addr := range addrs {
		if addr.To4() != nil {
			picked = addr
			break
		}
	}
	return picked.String(), ""
}

// ---------------------------------------------------------------- 数据源实现

// ip.sb —— 字段最全（经纬度 / 时区 / AS 编号都有），境外
type ipSBResponse struct {
	IP              string  `json:"ip"`
	Country         string  `json:"country"`
	CountryCode     string  `json:"country_code"`
	Region          string  `json:"region"`
	RegionCode      string  `json:"region_code"`
	City            string  `json:"city"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Timezone        string  `json:"timezone"`
	ISP             string  `json:"isp"`
	Organization    string  `json:"organization"`
	ASNOrganization string  `json:"asn_organization"`
	ASN             int     `json:"asn"`
}

func buildIPSBURL(query string) string {
	if query == "" {
		return "https://api.ip.sb/geoip"
	}
	return "https://api.ip.sb/geoip/" + url.PathEscape(query)
}

func parseIPSB(body []byte) (*IPInfo, error) {
	var raw ipSBResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	if raw.IP == "" {
		return nil, errors.New("响应缺少 ip 字段")
	}

	org := raw.ASNOrganization
	if org == "" {
		org = raw.Organization
	}
	return &IPInfo{
		Status:      "success",
		Country:     raw.Country,
		CountryCode: raw.CountryCode,
		Region:      raw.RegionCode,
		RegionName:  raw.Region,
		City:        raw.City,
		Lat:         raw.Latitude,
		Lon:         raw.Longitude,
		Timezone:    raw.Timezone,
		Isp:         raw.ISP,
		Org:         org,
		As:          formatASN(raw.ASN),
		Query:       raw.IP,
	}, nil
}

// vore.top —— 国内接口，返回中文地名，无经纬度 / 时区 / AS 编号
type voreResponse struct {
	Code   int `json:"code"`
	IPInfo struct {
		Text string `json:"text"`
	} `json:"ipinfo"`
	IPData struct {
		Info1 string `json:"info1"` // 国家
		Info2 string `json:"info2"` // 省 / 州
		Info3 string `json:"info3"` // 城市
		ISP   string `json:"isp"`
	} `json:"ipdata"`
}

func buildVoreURL(query string) string {
	if query == "" {
		return "https://api.vore.top/api/IPdata"
	}
	return "https://api.vore.top/api/IPdata?ip=" + url.QueryEscape(query)
}

func parseVore(body []byte) (*IPInfo, error) {
	var raw voreResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	if raw.Code != http.StatusOK || raw.IPData.Info1 == "" {
		return nil, fmt.Errorf("响应异常（code=%d）", raw.Code)
	}
	return &IPInfo{
		Status:     "success",
		Country:    raw.IPData.Info1,
		RegionName: raw.IPData.Info2,
		City:       raw.IPData.Info3,
		Isp:        raw.IPData.ISP,
		Query:      raw.IPInfo.Text,
	}, nil
}

// ipwho.is —— 字段较全，境外
type ipWhoResponse struct {
	Success     bool    `json:"success"`
	Message     string  `json:"message"`
	IP          string  `json:"ip"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Region      string  `json:"region"`
	RegionCode  string  `json:"region_code"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Connection  struct {
		ASN int    `json:"asn"`
		Org string `json:"org"`
		ISP string `json:"isp"`
	} `json:"connection"`
	Timezone struct {
		ID string `json:"id"`
	} `json:"timezone"`
}

func buildIPWhoURL(query string) string {
	if query == "" {
		return "https://ipwho.is/"
	}
	return "https://ipwho.is/" + url.PathEscape(query)
}

func parseIPWho(body []byte) (*IPInfo, error) {
	var raw ipWhoResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	if !raw.Success || raw.IP == "" {
		if raw.Message != "" {
			return nil, errors.New(raw.Message)
		}
		return nil, errors.New("响应异常")
	}

	isp := raw.Connection.ISP
	if isp == "" {
		isp = raw.Connection.Org
	}
	return &IPInfo{
		Status:      "success",
		Country:     raw.Country,
		CountryCode: raw.CountryCode,
		Region:      raw.RegionCode,
		RegionName:  raw.Region,
		City:        raw.City,
		Lat:         raw.Latitude,
		Lon:         raw.Longitude,
		Timezone:    raw.Timezone.ID,
		Isp:         isp,
		Org:         raw.Connection.Org,
		As:          formatASN(raw.Connection.ASN),
		Query:       raw.IP,
	}, nil
}

// formatASN 把 AS 编号格式化成 "AS15169"，无效编号返回空串
func formatASN(asn int) string {
	if asn <= 0 {
		return ""
	}
	return fmt.Sprintf("AS%d", asn)
}

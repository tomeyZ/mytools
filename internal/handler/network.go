package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

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

type NetworkHandler struct {
	client *http.Client
}

func NewNetworkHandler() *NetworkHandler {
	return &NetworkHandler{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetIpInfo 获取IP信息
func (a *NetworkHandler) GetIpInfo(ip string) (*IPInfo, error) {
	url := "http://ip-api.com/json/" + ip + "?lang=zh-CN"
	resp, err := a.client.Get(url)
	if err != nil {
		return &IPInfo{Status: "error", Message: err.Error()}, err
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return &IPInfo{Status: "error", Message: "请求查询失败"}, nil
	}

	var ipInfo IPInfo
	err = json.NewDecoder(resp.Body).Decode(&ipInfo)
	if err != nil {
		return &IPInfo{Status: "error", Message: err.Error()}, err
	}
	return &ipInfo, nil
}

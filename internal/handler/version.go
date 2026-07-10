package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mytools/internal/config"
	"net/http"
	"runtime"
	"time"
)

type ApiResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data VersionInfo `json:"data"`
}

type VersionInfo struct {
	ID         int      `json:"id"`
	Version    string   `json:"version"`
	ChangeLog  []string `json:"change_log"`
	CreateTime int64    `json:"create_time"`
	CreateDate string   `json:"create_date"`
}

type VersionHandler struct {
	client *http.Client
}

func NewVersionHandler() *VersionHandler {
	return &VersionHandler{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetVersionLatest 获取最新版本信息
func (a *VersionHandler) GetVersionLatest() (*VersionInfo, error) {
	url := "http://dabangke.com/api/v1/version/getLatest"
	resp, err := a.client.Post(url, "application/json", bytes.NewBuffer([]byte{}))
	//{"code":200,"msg":"success","data":{"id":1,"version":"1.0.2579","change_log":["新增JSON格式化功能","新增IP地址查询功能"],"create_time":1751881264,"create_date":"2025-07-07 17:41:04"}}
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request status code error: %v", resp.Status)
	}

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("json decode error: %v", err)
	}
	result := &apiResp.Data
	return result, nil
}

// CheckUpdate 检查更新
func (a *VersionHandler) CheckUpdate() (*VersionInfo, error) {
	//time.Sleep(10 * time.Second)
	versionInfo, err := a.GetVersionLatest()
	if err != nil {
		return nil, err
	}
	if versionInfo == nil {
		return nil, nil
	}
	currentVersion := config.Get().Version
	if currentVersion == versionInfo.Version {
		return nil, nil
	}
	return versionInfo, nil
}

// Download 下载最新版本
func (a *VersionHandler) Download() ([]byte, error) {
	cfg := config.Get()
	extension := GetExtension()
	versionInfo, err := a.CheckUpdate()
	if err != nil || versionInfo == nil || extension == "" {
		return nil, err
	}

	downloadUrl := cfg.BaseUrl + "/downloads/" + versionInfo.Version + "/MyTools_v" + versionInfo.Version + extension
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func GetExtension() string {
	switch runtime.GOOS {
	case "windows":
		return ".exe"
	case "darwin":
		return ".dmg"
	default:
		return ""
	}
}

// GetCurrentVersion 获取当前版本号
func (a *VersionHandler) GetCurrentVersion() string {
	return config.Get().Version
}

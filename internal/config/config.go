package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type AuthorConfig struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AppConfig struct {
	Name    string       `json:"name"`
	Version string       `json:"version"`
	BaseUrl string       `json:"baseUrl"`
	Author  AuthorConfig `json:"author"`
}

var (
	instance *AppConfig
	once     sync.Once
	initErr  error
)

// 配置源是根目录的 wails.json：main.go 里 //go:embed 嵌入后传入字节。
// 不在本包再放一份副本——go:embed 只能嵌包内文件，副本会导致版本号双写漂移
func Init(configData []byte) error {
	once.Do(func() {
		instance, initErr = loadConfig(configData)
	})
	return initErr
}

func loadConfig(configData []byte) (*AppConfig, error) {
	var config AppConfig

	// 解析配置
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("配置解析失败：%v", err)
	}
	log.Println("配置加载成功！")
	return &config, nil
}

func Get() *AppConfig {
	if instance == nil {
		log.Fatal("配置初始化失败")
	}
	return instance
}

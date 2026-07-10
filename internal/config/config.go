package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

//go:embed wails.json
var embeddedConfig embed.FS

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

func LoadConfig() (*AppConfig, error) {
	var config AppConfig

	// 从嵌入的文件系统中读取 wails.json
	configData, err := embeddedConfig.ReadFile("wails.json")
	if err != nil {
		return nil, fmt.Errorf("读取内嵌配置文件失败：%v", err)
	}

	// 解析配置
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("配置解析失败：%v", err)
	}
	log.Println("配置加载成功！")
	return &config, nil
}

func Init() error {
	once.Do(func() {
		instance, initErr = LoadConfig()
	})
	return initErr
}

func Get() *AppConfig {
	if instance == nil {
		log.Fatal("配置初始化失败")
	}
	return instance
}

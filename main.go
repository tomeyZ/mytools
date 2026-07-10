package main

import (
	"embed"
	"log"
	"mytools/internal/config"
	"mytools/internal/handler"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatal("配置初始化失败: ", err)
	}
	cfg := config.Get()

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     cfg.Name,
		Width:     1024,
		Height:    768,
		MinWidth:  1024,
		MinHeight: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			handler.NewNetworkHandler(),
			handler.NewTimeHandler(),
			handler.NewVersionHandler(),
			handler.NewCryptoHandler(),
			handler.NewRsaHandler(),
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

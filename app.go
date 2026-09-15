package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// OpenExternal 用系统默认浏览器打开外部链接（下载新版、Release 页面）。
// 不用前端 window.open：WebView2 里会在应用窗口内打开，没有地址栏也没有下载能力
func (a *App) OpenExternal(url string) {
	if a.ctx == nil || url == "" {
		return
	}
	runtime.BrowserOpenURL(a.ctx, url)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

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

// SaveQRCode 把前端生成的二维码（data URL）写到用户选定的位置。
func (a *App) SaveQRCode(dataURL string, defaultName string) (string, error) {
	const prefix = "data:image/png;base64,"

	if a.ctx == nil {
		return "", fmt.Errorf("应用尚未就绪")
	}
	if !strings.HasPrefix(dataURL, prefix) {
		return "", fmt.Errorf("图片数据格式不正确")
	}
	raw, err := base64.StdEncoding.DecodeString(dataURL[len(prefix):])
	if err != nil {
		return "", fmt.Errorf("图片数据解析失败: %w", err)
	}

	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "保存二维码",
		DefaultFilename: defaultName,
		Filters: []runtime.FileFilter{
			{DisplayName: "PNG 图片", Pattern: "*.png"},
		},
	})
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", nil
	}
	// 用户可能在对话框里把扩展名删掉或改成别的，统一补回 .png
	if !strings.HasSuffix(strings.ToLower(path), ".png") {
		path += ".png"
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}
	return path, nil
}

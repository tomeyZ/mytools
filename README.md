# 开发工具（mytools）

一款基于 Wails v2 + Vue 3 的桌面端开发工具箱，集成了日常开发中常用的小工具，界面简洁、开箱即用。

## 功能特性

- **时区转换** —— 时间戳与日期互转，支持 8 个常用时区、日期计算
- **MD5 加密** —— 文本 MD5 摘要计算
- **AES 加解密** —— 支持 CBC / ECB 模式，PKCS7 填充
- **RSA 密钥生成** —— 生成 RSA 密钥对，支持 PKCS#1 / PKCS#8
- **JSON 美化** —— 格式化、压缩与校验
- **文本处理** —— 常用文本转换与处理
- **二维码** —— 识别图片中的二维码，或由文本生成二维码并导出 PNG
- **IP 地址查询** —— 本机 IP 与归属地查询
- **应用内更新** —— 自动检测 GitHub Releases 最新版本并提示下载

界面支持浅色 / 深色主题切换，`Ctrl K` 全局搜索快速定位工具。

## 下载安装

前往 [Releases](https://github.com/tomeyZ/mytools/releases) 页面下载对应平台的安装包：

| 平台 | 文件 | 说明 |
|------|------|------|
| Windows | `mytools_v*.exe` | 单文件，下载后直接运行，无需安装 |
| macOS | `mytools_*.dmg` | Apple Silicon / Intel 通用 |

> Windows 10 若提示无法启动，请先安装 [WebView2 Runtime](https://developer.microsoft.com/microsoft-edge/webview2/)（Windows 11 自带）。

应用启动时会静默检查一次更新。有新版本时，侧边栏左下角的「检查更新」会挂上红点，点击即可查看更新日志并跳转下载。

## 本地开发

### 环境要求

- [Go](https://go.dev/) 1.25+
- [Node.js](https://nodejs.org/) 18+
- [Wails CLI v2](https://wails.io/)

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 常用命令

```bash
# 开发模式（热重载）
wails dev

# 构建当前平台
wails build

# Windows 构建 NSIS 安装包（可选）
wails build -nsis
```

## 发布流程

推送 `v` 开头的 tag 即触发 CI，自动构建两个平台的产物并创建 Release。

```bash
# 1. 把 wails.json 的 version 改成新版本号（不带 v）并提交
git commit -am "chore: 版本号 1.1.5"

# 2. 打 tag 并推送，需与 wails.json 保持一致
git tag v1.1.5
git push origin main --tags
```

Release 正文会作为更新日志显示在应用的更新提示里。

## 项目结构

```
mytools/
├── main.go                 # 应用入口，前端资源与配置内嵌
├── app.go                  # 应用生命周期
├── wails.json              # 应用名与版本号
├── internal/
│   ├── config/             # 配置加载（wails.json 编译期内嵌）
│   └── handler/            # 绑定给前端的后端方法
│       ├── crypto.go       # AES
│       ├── rsa.go          # RSA
│       ├── time.go         # 时间转换
│       ├── network.go      # IP 查询
│       └── version.go      # 版本检查与更新
├── frontend/               # Vue 3 + Vite 前端
│   └── src/
│       ├── tools.js        # 工具注册表（侧栏菜单与组件映射）
│       ├── views/          # 整体布局
│       └── components/     # 各工具页面组件
├── design/                 # UI 设计稿（不参与构建）
└── .github/workflows/      # CI 自动构建
```

## License

[MIT](LICENSE)

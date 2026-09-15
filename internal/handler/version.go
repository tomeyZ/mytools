package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"mytools/internal/config"
)

const (
	repoOwner = "tomeyZ"
	repoName  = "mytools"

	// GitHub Releases 最新公开版本（不含 draft / prerelease）
	releasesAPI = "https://api.github.com/repos/" + repoOwner + "/" + repoName + "/releases/latest"

	// 兜底下载地址模板：CI 产物名是固定的，可以直接拼，
	// 省掉一次 assets 列表请求。见 .github/workflows/build-release.yml
	repoReleases = "https://github.com/" + repoOwner + "/" + repoName + "/releases/download/"
)

// VersionInfo 最新版本信息
type VersionInfo struct {
	Version     string   `json:"version"`
	ChangeLog   []string `json:"change_log"`
	CreateDate  string   `json:"create_date"`
	DownloadURL string   `json:"download_url"`
	ReleaseURL  string   `json:"release_url"`
}

// UpdateResult 检查结果。status 取值：
//
//	update —— 有新版本，Latest 里有版本、更新日志、下载地址
//	latest —— 已是最新
//	none   —— 仓库还没有 publish 过的 release（没有可下载的安装包）
//	error  —— 网络不通 / 接口异常，Message 是给用户看的中文
//
// 网络类失败一律走 status + 中文 Message，不要回传 error：
// Wails 会把 error 抛给前端 reject，用户看到的是 `Get "https://...": dial tcp` 这种英文原文。
type UpdateResult struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Latest  *VersionInfo `json:"latest"`
}

// httpStatusError 让调用方能分辨 HTTP 状态码（404 = 尚无 release，不是网络故障）
type httpStatusError struct {
	Code int
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("github api status: %d", e.Code)
}

// ghRelease 只取用得上的字段，其余交给 encoding/json 忽略
type ghRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	HTMLURL     string    `json:"html_url"`
	PublishedAt string    `json:"published_at"`
	Assets      []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type VersionHandler struct {
	client *http.Client
}

func NewVersionHandler() *VersionHandler {
	return &VersionHandler{
		client: &http.Client{Timeout: 8 * time.Second},
	}
}

// GetCurrentVersion 获取当前版本号
func (a *VersionHandler) GetCurrentVersion() string {
	return config.Get().Version
}

// CheckUpdate 检查是否有新版本。永不返回 error，失败降级为 status="error"
func (a *VersionHandler) CheckUpdate() UpdateResult {
	rel, err := a.fetchLatest()
	if err != nil {
		// 仓库有 tag 但没 publish 过 release 时，/releases/latest 返回 404 —— 这不是故障，
		// 当成"没有可升级版本"处理，别让用户在界面上看到"网络失败"
		var se *httpStatusError
		if errors.As(err, &se) && se.Code == http.StatusNotFound {
			return UpdateResult{Status: "none", Message: "暂未发布可用版本"}
		}
		return UpdateResult{Status: "error", Message: "检查更新失败，请检查网络后重试"}
	}

	latest := normalizeVersion(rel.TagName)
	if compareVersion(latest, config.Get().Version) <= 0 {
		return UpdateResult{Status: "latest", Message: "已是最新版本"}
	}

	return UpdateResult{
		Status:  "update",
		Message: "发现新版本",
		Latest: &VersionInfo{
			Version:     latest,
			ChangeLog:   parseChangeLog(rel.Body),
			CreateDate:  formatDate(rel.PublishedAt),
			DownloadURL: pickAssetURL(rel.Assets, rel.TagName),
			ReleaseURL:  rel.HTMLURL,
		},
	}
}

func (a *VersionHandler) fetchLatest() (*ghRelease, error) {
	req, err := http.NewRequest(http.MethodGet, releasesAPI, nil)
	if err != nil {
		return nil, err
	}
	// GitHub API 不带 User-Agent 会直接 403，且必须声明版本才返回稳定字段
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", config.Get().Name)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &httpStatusError{Code: resp.StatusCode}
	}

	var rel ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}
	if rel.TagName == "" {
		return nil, fmt.Errorf("github api 返回空 tag")
	}
	return &rel, nil
}

// normalizeVersion 去掉 tag 前缀（v1.1.4 -> 1.1.4）与空白
func normalizeVersion(tag string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(tag), "v"))
}

// compareVersion 语义化版本号比较：a>b 返回 1，a<b 返回 -1，相等返回 0。
// 非数字段退化为字符串比较（如 1.1.4-beta）
func compareVersion(a, b string) int {
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	n := len(as)
	if len(bs) > n {
		n = len(bs)
	}
	for i := 0; i < n; i++ {
		var x, y string
		if i < len(as) {
			x = as[i]
		}
		if i < len(bs) {
			y = bs[i]
		}
		xi, xe := strconv.Atoi(x)
		yi, ye := strconv.Atoi(y)
		if xe == nil && ye == nil {
			if xi != yi {
				if xi > yi {
					return 1
				}
				return -1
			}
			continue
		}
		if x != y {
			if x > y {
				return 1
			}
			return -1
		}
	}
	return 0
}

// pickAssetURL 挑出当前平台的安装包地址；找不到返回空串，前端会退回 Release 页面
func pickAssetURL(assets []ghAsset, tag string) string {
	var keyword, suffix string
	switch runtime.GOOS {
	case "windows":
		keyword, suffix = "windows", ".exe"
	case "darwin":
		keyword, suffix = "darwin", ".dmg"
	default:
		return ""
	}

	for _, it := range assets {
		name := strings.ToLower(it.Name)
		if strings.Contains(name, keyword) && strings.HasSuffix(name, suffix) {
			if it.BrowserDownloadURL != "" {
				return it.BrowserDownloadURL
			}
			return repoReleases + tag + "/" + it.Name
		}
	}

	// 兜底：按 CI 的命名模板直接拼
	return repoReleases + tag + "/" + assetFileName(tag)
}

// assetFileName CI 产物名模板：mytools_v1.1.4_windows_amd64.exe / _darwin_universal.dmg
func assetFileName(tag string) string {
	switch runtime.GOOS {
	case "windows":
		arch := runtime.GOARCH
		if arch != "amd64" && arch != "386" {
			arch = "amd64"
		}
		return fmt.Sprintf("mytools_%s_windows_%s.exe", tag, arch)
	case "darwin":
		return fmt.Sprintf("mytools_%s_darwin_universal.dmg", tag)
	default:
		return ""
	}
}

// parseChangeLog 把 Release 正文拆成更新列表：去空行、去列表符号、最多 20 条
func parseChangeLog(body string) []string {
	var out []string
	for _, line := range strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n") {
		s := strings.TrimSpace(line)
		s = strings.TrimSpace(strings.TrimPrefix(s, "-"))
		s = strings.TrimSpace(strings.TrimPrefix(s, "*"))
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		out = append(out, s)
		if len(out) >= 20 {
			break
		}
	}
	return out
}

// formatDate 2026-09-15T08:12:33Z -> 2026-09-15
func formatDate(raw string) string {
	if len(raw) < 10 {
		return ""
	}
	return raw[:10]
}

package handler

import (
	"fmt"
	"time"
)

type TimeHandler struct {
}

func NewTimeHandler() *TimeHandler {
	return &TimeHandler{}
}

// TimestampToDate 时间戳转日期（支持秒/毫秒，自动判别）
func (a *TimeHandler) TimestampToDate(timestamp int64, timezone string) string {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "无效时区"
	}
	// 13位及以上按毫秒处理，10位按秒
	sec := timestamp
	if timestamp > 1e12 || timestamp < -1e12 {
		sec = timestamp / 1000
	}
	return time.Unix(sec, 0).In(loc).Format(time.DateTime)
}

// DateToTimestamp 日期转时间戳，失败返回错误信息
func (a *TimeHandler) DateToTimestamp(dateStr, timezone string) string {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "无效时区"
	}
	t, err := time.ParseInLocation(time.DateTime, dateStr, loc)
	if err != nil {
		return "格式错误，应为 YYYY-MM-DD HH:mm:ss"
	}
	return fmt.Sprintf("%d", t.Unix())
}

// GetCurrentTime 获取当前时间
func (a *TimeHandler) GetCurrentTime(timezone string) map[string]interface{} {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return map[string]interface{}{
			"error": "无效时区",
		}
	}
	now := time.Now().In(loc)
	return map[string]interface{}{
		"datetime":  now.Format(time.DateTime),
		"timestamp": now.Unix(),
	}
}

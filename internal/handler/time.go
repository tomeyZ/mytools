package handler

import "time"

type TimeHandler struct {
}

func NewTimeHandler() *TimeHandler {
	return &TimeHandler{}
}

// TimestampToDate 时间戳转日期
func (a *TimeHandler) TimestampToDate(timestamp int64, timezone string) string {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "无效时区"
	}
	return time.Unix(timestamp, 0).In(loc).Format(time.DateTime)
}

// DateToTimestamp 日期转时间戳
func (a *TimeHandler) DateToTimestamp(dateStr, timezone string) int64 {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return -1
	}
	t, err := time.ParseInLocation(time.DateTime, dateStr, loc)
	if err != nil {
		return -1
	}
	return t.Unix()
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

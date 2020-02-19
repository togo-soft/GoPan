package utils

import (
	"time"
)

// GetNowDateTime 返回当前日期时间
func GetNowDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// Unix2DateTime 时间戳转日期时间
func Unix2DateTime(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

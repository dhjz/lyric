package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return strconv.FormatUint(bytes, 10) + " B"
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func FormatDuration(duration time.Duration) string {
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	return fmt.Sprintf("%d天%d小时%d分钟", days, hours, minutes)
}

// 传入多个可能的参数, 返回一个string
func GetParam(r *http.Request, params ...string) string {
	for _, param := range params {
		if value := r.URL.Query().Get(param); value != "" {
			return value
		}
	}
	return ""
}

func GetParamInt(r *http.Request, param string, defaultVal int) int {
	intValue, err := strconv.Atoi(r.URL.Query().Get(param))
	if err == nil {
		return intValue
	}
	return defaultVal
}

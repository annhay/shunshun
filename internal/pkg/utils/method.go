package utils

import "time"

// StringTransformationTime 字符串转时间
func StringTransformationTime(str string) time.Time {
	// 解析时间字符串，指定本地时区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	times, err := time.ParseInLocation("2006-01-02 15:04:05", str, loc)
	if err != nil {
		return time.Time{}
	}
	return times
}

// TimeTransformationString 时间转字符串
func TimeTransformationString(now time.Time) string {
	return now.Format("2006-01-02 15:04:05")
}

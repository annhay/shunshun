package utils

import (
	"fmt"
	"math/rand"
	"time"
)

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

// DriverNoRandom 司机编号生成（唯一性）
func DriverNoRandom(userId int64, cityCode string) string {
	timestamp := time.Now().UnixNano()
	random := rand.Intn(10000)
	driverNo := fmt.Sprintf("Dri%d_%s_%d_%04d", userId, cityCode, timestamp, random)
	return driverNo
}

// OrderCodeRandom 订单编号随机生成（唯一性）
func OrderCodeRandom(userId int64) string {
	timestamp := time.Now().UnixNano()
	random := rand.Intn(10000)
	orderCode := fmt.Sprintf("SS%d_%d_%04d", userId, timestamp, random)
	return orderCode
}

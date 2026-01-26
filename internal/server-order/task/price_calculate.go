package task

import (
	"shunshun/internal/pkg/utils"
	"strconv"
)

// CalculateTravelDistance 计算行程距离
func CalculateTravelDistance(startLng, startLat, endLng, endLat float64) (float64, error) {
	if startLng == 0 || startLat == 0 || endLng == 0 || endLat == 0 {
		return 0, nil
	}
	
	// 转换为字符串格式
	startLngStr := strconv.FormatFloat(startLng, 'f', 6, 64)
	startLatStr := strconv.FormatFloat(startLat, 'f', 6, 64)
	endLngStr := strconv.FormatFloat(endLng, 'f', 6, 64)
	endLatStr := strconv.FormatFloat(endLat, 'f', 6, 64)
	
	// 调用地图服务计算距离
	drivingResult, err := utils.DrivingRoute(startLngStr, startLatStr, endLngStr, endLatStr)
	if err != nil {
		return 0, err
	}
	
	// 转换为公里
	distance, _ := strconv.ParseFloat(drivingResult.Distance, 64)
	travelDistance := distance / 1000
	
	return travelDistance, nil
}

// CalculateEstimatedAmount 计算预估金额
func CalculateEstimatedAmount(travelDistance float64, carType string) float64 {
	var basePrice float64
	var distancePrice float64
	
	// 根据车辆类型调整价格
	switch carType {
	case "1": // 经济车
		basePrice = 10.0
		distancePrice = 2.5
	case "2": // 商务车
		basePrice = 15.0
		distancePrice = 3.5
	case "3": // 六座车
		basePrice = 20.0
		distancePrice = 4.0
	default:
		basePrice = 10.0
		distancePrice = 2.5
	}
	
	// 计算预估金额
	var estimatedAmount float64
	if travelDistance <= 3 {
		estimatedAmount = basePrice
	} else {
		estimatedAmount = basePrice + (travelDistance-3)*distancePrice
	}
	
	return estimatedAmount
}

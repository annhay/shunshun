// Package utils 包含交通相关的工具函数
package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"shunshun/internal/pkg/global"
)

// -------------------------- 路况分析相关 --------------------------

// TrafficConditionResponse 路况分析响应结构
// 这里假设使用高德地图的路况API
// 实际API可能需要根据具体地图服务提供商进行调整
type TrafficConditionResponse struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Roads    []struct {
		Name        string `json:"name"`        // 道路名称
		RoadStatus  string `json:"road_status"`  // 路况状态：0畅通，1缓行，2拥堵，3严重拥堵
		Speed       string `json:"speed"`        // 平均车速
		JamLength   string `json:"jam_length"`   // 拥堵长度
		Direction   string `json:"direction"`    // 方向
		Polyline    string `json:"polyline"`     // 道路几何信息
	} `json:"roads"`
}

// TrafficConditionResult 路况分析结果
type TrafficConditionResult struct {
	Roads       []RoadInfo `json:"roads"`        // 道路信息
	AverageSpeed float64    `json:"average_speed"` // 平均车速
	CongestionLevel int     `json:"congestion_level"` // 整体拥堵等级：0畅通，1缓行，2拥堵，3严重拥堵
}

// RoadInfo 道路信息
type RoadInfo struct {
	Name        string  `json:"name"`
	RoadStatus  int     `json:"road_status"`
	Speed       float64 `json:"speed"`
	JamLength   float64 `json:"jam_length"`
	Direction   string  `json:"direction"`
}

// -------------------------- API URL构建函数 --------------------------

// buildTrafficURL 构建路况分析API URL
// 路况分析：获取指定路线的实时路况
func buildTrafficURL() string {
	return getAmapBaseURL() + "/v3/traffic/status"
}

// -------------------------- 路况分析相关函数 --------------------------

// GetTrafficCondition 获取指定路线的实时路况
func GetTrafficCondition(startLng, startLat, endLng, endLat float64) (*TrafficConditionResult, error) {
	// 生成缓存键
	cacheKey := fmt.Sprintf("traffic:%f,%f,%f,%f", startLng, startLat, endLng, endLat)
	cacheExpiration := 5 * time.Minute // 路况缓存5分钟

	// 检查缓存
	if global.Rdb != nil {
		cachedData, err := global.Rdb.Get(global.Ctx, cacheKey).Result()
		if err == nil {
			// 缓存命中
			var result TrafficConditionResult
			if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
				return &result, nil
			}
		}
	}

	// 构建请求参数
	start := fmt.Sprintf("%f,%f", startLng, startLat)
	end := fmt.Sprintf("%f,%f", endLng, endLat)
	encodedStart := url.QueryEscape(start)
	encodedEnd := url.QueryEscape(end)

	// 构建请求URL
	reqURL := fmt.Sprintf("%s?key=%s&extensions=all&level=1&origin=%s&destination=%s", 
		buildTrafficURL(), getAmapAPIKey(), encodedStart, encodedEnd)

	// 发送请求
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析响应
	body, _ := ioutil.ReadAll(resp.Body)
	var res TrafficConditionResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("解析路况响应失败: %v", err)
	}

	if res.Status != "1" {
		return nil, fmt.Errorf("获取路况失败: %s", res.Info)
	}

	// 处理结果
	var roads []RoadInfo
	totalSpeed := 0.0
	maxCongestion := 0

	for _, road := range res.Roads {
		roadStatus, _ := strconv.Atoi(road.RoadStatus)
		speed, _ := strconv.ParseFloat(road.Speed, 64)
		jamLength, _ := strconv.ParseFloat(road.JamLength, 64)

		roads = append(roads, RoadInfo{
			Name:        road.Name,
			RoadStatus:  roadStatus,
			Speed:       speed,
			JamLength:   jamLength,
			Direction:   road.Direction,
		})

		totalSpeed += speed
		if roadStatus > maxCongestion {
			maxCongestion = roadStatus
		}
	}

	// 计算平均车速
	averageSpeed := 0.0
	if len(roads) > 0 {
		averageSpeed = totalSpeed / float64(len(roads))
	}

	result := &TrafficConditionResult{
		Roads:          roads,
		AverageSpeed:   averageSpeed,
		CongestionLevel: maxCongestion,
	}

	// 将结果存入缓存
	if global.Rdb != nil {
		data, err := json.Marshal(result)
		if err == nil {
			global.Rdb.Set(global.Ctx, cacheKey, data, cacheExpiration)
		}
	}

	return result, nil
}

// CalculateDynamicPrice 根据路况动态调整价格
func CalculateDynamicPrice(basePrice float64, trafficCondition *TrafficConditionResult) float64 {
	if trafficCondition == nil {
		return basePrice
	}

	// 根据拥堵等级调整价格
	// 畅通：原价
	// 缓行：加价10%
	// 拥堵：加价20%
	// 严重拥堵：加价30%
	var priceAdjustment float64
	switch trafficCondition.CongestionLevel {
	case 0: // 畅通
		priceAdjustment = 1.0
	case 1: // 缓行
		priceAdjustment = 1.1
	case 2: // 拥堵
		priceAdjustment = 1.2
	case 3: // 严重拥堵
		priceAdjustment = 1.3
	default:
		priceAdjustment = 1.0
	}

	// 应用价格调整
	dynamicPrice := basePrice * priceAdjustment

	return dynamicPrice
}

// GetAlternativeRoute 获取备选路线
func GetAlternativeRoute(startLng, startLat, endLng, endLat float64) ([]RouteOption, error) {
	// 这里简化实现，实际应该调用地图服务的路线规划API
	// 获取多条备选路线并根据路况、距离等因素排序
	
	// 转换为字符串格式
	startLngStr := strconv.FormatFloat(startLng, 'f', 6, 64)
	startLatStr := strconv.FormatFloat(startLat, 'f', 6, 64)
	endLngStr := strconv.FormatFloat(endLng, 'f', 6, 64)
	endLatStr := strconv.FormatFloat(endLat, 'f', 6, 64)
	
	// 调用地图服务获取备选路线
	drivingResult, err := DrivingRoute(startLngStr, startLatStr, endLngStr, endLatStr)
	if err != nil {
		return nil, err
	}

	// 这里假设只返回一条路线
	// 实际实现中应该获取多条路线并进行比较
	distance, _ := strconv.ParseFloat(drivingResult.Distance, 64)
	duration, _ := strconv.ParseFloat(drivingResult.Duration, 64)

	// 获取路况
	trafficCondition, err := GetTrafficCondition(startLng, startLat, endLng, endLat)
	if err != nil {
		// 路况获取失败，使用默认值
		trafficCondition = &TrafficConditionResult{
			CongestionLevel: 0,
			AverageSpeed:    60.0,
		}
	}

	// 计算基础价格
	basePrice := CalculateEstimatedAmount(distance/1000, "1") // 假设使用经济车
	
	// 动态调整价格
	dynamicPrice := CalculateDynamicPrice(basePrice, trafficCondition)

	// 构建备选路线
	routeOptions := []RouteOption{
		{
			Distance:        distance / 1000, // 转换为公里
			Duration:        duration / 60,    // 转换为分钟
			Price:           dynamicPrice,
			CongestionLevel: trafficCondition.CongestionLevel,
			AverageSpeed:    trafficCondition.AverageSpeed,
			RouteType:       "fastest",
			Description:     "最快路线",
		},
	}

	return routeOptions, nil
}

// RouteOption 路线选项
type RouteOption struct {
	Distance        float64 `json:"distance"`        // 距离（公里）
	Duration        float64 `json:"duration"`        // 时长（分钟）
	Price           float64 `json:"price"`           // 预估价格
	CongestionLevel int     `json:"congestion_level"` // 拥堵等级
	AverageSpeed    float64 `json:"average_speed"`    // 平均车速
	RouteType       string  `json:"route_type"`       // 路线类型：fastest, shortest, economical
	Description     string  `json:"description"`      // 路线描述
}

// CalculateEstimatedAmount 计算预估金额
// 注意：这里复用了原有的金额计算函数
// 实际实现中可能需要根据具体业务逻辑进行调整
func CalculateEstimatedAmount(travelDistance float64, carType string) float64 {
	// 这里可以根据需要实现更复杂的计价逻辑
	// 暂时返回一个基础价格
	var basePrice float64
	var distancePrice float64
	
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
	
	var estimatedAmount float64
	if travelDistance <= 3 {
		estimatedAmount = basePrice
	} else {
		estimatedAmount = basePrice + (travelDistance-3)*distancePrice
	}
	
	return estimatedAmount
}

// Package utils 包含AI天气提醒服务，使用通义千问生成智能出行建议
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"shunshun/internal/pkg/global"
)

// -------------------------- 通义千问AI配置 --------------------------

// getTongyiAPIKey 获取通义千问API密钥
func getTongyiAPIKey() string {
	if global.AppConf != nil {
		return global.AppConf.Tongyi.APIKey
	}
	return ""
}

// getTongyiBaseURL 获取通义千问API基础URL
func getTongyiBaseURL() string {
	if global.AppConf != nil && global.AppConf.Tongyi.BaseURL != "" {
		return global.AppConf.Tongyi.BaseURL
	}
	return "https://dashscope.aliyuncs.com/api/v1"
}

// getTongyiModel 获取通义千问模型名称
func getTongyiModel() string {
	if global.AppConf != nil && global.AppConf.Tongyi.Model != "" {
		return global.AppConf.Tongyi.Model
	}
	return "qwen-turbo"
}

// -------------------------- AI请求/响应结构 --------------------------

// TongyiRequest 通义千问请求结构
type TongyiRequest struct {
	Model string `json:"model"`
	Input struct {
		Messages []TongyiMessage `json:"messages"`
	} `json:"input"`
	Parameters struct {
		ResultFormat string `json:"result_format"`
	} `json:"parameters"`
}

// TongyiMessage 消息结构
type TongyiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// TongyiResponse 通义千问响应结构
type TongyiResponse struct {
	Output struct {
		Text         string `json:"text"`
		FinishReason string `json:"finish_reason"`
		Choices      []struct {
			FinishReason string `json:"finish_reason"`
			Message      struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	} `json:"output"`
	Usage struct {
		OutputTokens int `json:"output_tokens"`
		InputTokens  int `json:"input_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
	RequestID string `json:"request_id"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

// -------------------------- AI天气提醒服务 --------------------------

// WeatherReminder AI天气提醒结果
type WeatherReminder struct {
	City          string   `json:"city"`          // 城市
	Weather       string   `json:"weather"`       // 天气现象
	Temperature   string   `json:"temperature"`   // 温度
	Humidity      string   `json:"humidity"`      // 湿度
	Winddirection string   `json:"winddirection"` // 风向
	Windpower     string   `json:"windpower"`     // 风力
	AIAdvice      string   `json:"ai_advice"`     // AI出行建议
	SafetyTips    []string `json:"safety_tips"`   // 安全提示
	DrivingLevel  string   `json:"driving_level"` // 驾驶风险等级
	UpdateTime    string   `json:"update_time"`   // 更新时间
}

// GetWeatherReminder 获取AI智能天气提醒
// adcode: 城市行政编码
func GetWeatherReminder(adcode string) (*WeatherReminder, error) {
	// 1. 首先获取实时天气数据
	weather, err := WeatherQuery(adcode)
	if err != nil {
		return nil, fmt.Errorf("获取天气数据失败: %v", err)
	}

	// 2. 调用AI生成出行建议
	aiAdvice, err := generateAIAdvice(weather)
	if err != nil {
		// AI服务不可用时，使用规则生成建议
		aiAdvice = generateRuleBasedAdvice(weather)
	}

	// 3. 生成安全提示和驾驶风险等级
	safetyTips, drivingLevel := generateSafetyInfo(weather)

	return &WeatherReminder{
		City:          weather.City,
		Weather:       weather.Weather,
		Temperature:   weather.Temperature,
		Humidity:      weather.Humidity,
		Winddirection: weather.Winddirection,
		Windpower:     weather.Windpower,
		AIAdvice:      aiAdvice,
		SafetyTips:    safetyTips,
		DrivingLevel:  drivingLevel,
		UpdateTime:    weather.Reporttime,
	}, nil
}

// GetWeatherReminderByAddress 根据地址获取AI智能天气提醒
func GetWeatherReminderByAddress(address string) (*WeatherReminder, error) {
	// 1. 地理编码获取城市编码
	geoResult, err := Geocode(address)
	if err != nil {
		return nil, fmt.Errorf("地址解析失败: %v", err)
	}

	// 2. 获取天气提醒
	return GetWeatherReminder(geoResult.Adcode)
}

// GetWeatherReminderByLocation 根据经纬度获取AI智能天气提醒
func GetWeatherReminderByLocation(lng, lat string) (*WeatherReminder, error) {
	// 1. 逆地理编码获取城市编码
	regeoResult, err := Regeocode(lng, lat)
	if err != nil {
		return nil, fmt.Errorf("位置解析失败: %v", err)
	}

	// 2. 获取天气提醒
	return GetWeatherReminder(regeoResult.Adcode)
}

// generateAIAdvice 调用通义千问生成AI出行建议
func generateAIAdvice(weather *WeatherResult) (string, error) {
	apiKey := getTongyiAPIKey()
	if apiKey == "" || apiKey == "sk-xxxxx" {
		return "", fmt.Errorf("API Key未配置")
	}

	// 构建提示词
	prompt := fmt.Sprintf(`你是一位专业的网约车出行顾问。请根据以下天气信息，为网约车乘客和司机提供简洁实用的出行建议（100字以内）：

城市：%s
天气：%s
温度：%s℃
湿度：%s%%
风向：%s
风力：%s级

请从以下几个方面给出建议：
1. 是否适合出行
2. 乘车注意事项
3. 司机驾驶建议
4. 特殊天气应对

请用简洁的中文回答，语气温馨友好。`,
		weather.City, weather.Weather, weather.Temperature,
		weather.Humidity, weather.Winddirection, weather.Windpower)

	// 构建请求
	reqBody := TongyiRequest{
		Model: getTongyiModel(),
	}
	reqBody.Input.Messages = []TongyiMessage{
		{Role: "system", Content: "你是一位专业的网约车出行顾问，擅长根据天气情况给出安全实用的出行建议。"},
		{Role: "user", Content: prompt},
	}
	reqBody.Parameters.ResultFormat = "message"

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// 发送请求
	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", getTongyiBaseURL()+"/services/aigc/text-generation/generation", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result TongyiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析AI响应失败: %v", err)
	}

	if result.Code != "" {
		return "", fmt.Errorf("AI服务错误: %s - %s", result.Code, result.Message)
	}

	// 提取AI回复
	if len(result.Output.Choices) > 0 {
		return result.Output.Choices[0].Message.Content, nil
	}
	if result.Output.Text != "" {
		return result.Output.Text, nil
	}

	return "", fmt.Errorf("AI未返回有效内容")
}

// generateRuleBasedAdvice 基于规则生成出行建议（AI服务不可用时的备选方案）
func generateRuleBasedAdvice(weather *WeatherResult) string {
	weatherType := weather.Weather

	switch {
	case contains(weatherType, "雨"):
		return fmt.Sprintf("【%s天气提醒】当前%s，气温%s℃。雨天路滑，建议乘客提前预约用车，司机朋友请减速慢行，保持安全车距。注意携带雨具，祝您出行顺利！",
			weather.City, weather.Weather, weather.Temperature)
	case contains(weatherType, "雪"):
		return fmt.Sprintf("【%s天气提醒】当前%s，气温%s℃。雪天道路湿滑，建议非必要不出行。如需用车，请预留充足时间，司机请谨慎驾驶，避免急刹车急转弯。",
			weather.City, weather.Weather, weather.Temperature)
	case contains(weatherType, "雾") || contains(weatherType, "霾"):
		return fmt.Sprintf("【%s天气提醒】当前%s，能见度较低。建议出行佩戴口罩，司机请开启雾灯，保持低速行驶。如遇严重雾霾，建议延迟出行。",
			weather.City, weather.Weather)
	case contains(weatherType, "晴"):
		return fmt.Sprintf("【%s天气提醒】当前%s，气温%s℃，非常适合出行！温馨提示：注意防晒补水，司机朋友注意避免阳光直射影响视线。祝您旅途愉快！",
			weather.City, weather.Weather, weather.Temperature)
	case contains(weatherType, "多云") || contains(weatherType, "阴"):
		return fmt.Sprintf("【%s天气提醒】当前%s，气温%s℃，适合出行。建议随身携带雨具以备不时之需。祝您出行顺利！",
			weather.City, weather.Weather, weather.Temperature)
	default:
		return fmt.Sprintf("【%s天气提醒】当前%s，气温%s℃，湿度%s%%。请根据天气情况合理安排出行，祝您一路平安！",
			weather.City, weather.Weather, weather.Temperature, weather.Humidity)
	}
}

// generateSafetyInfo 生成安全提示和驾驶风险等级
func generateSafetyInfo(weather *WeatherResult) ([]string, string) {
	var tips []string
	var level string
	weatherType := weather.Weather

	switch {
	case contains(weatherType, "暴雨") || contains(weatherType, "暴雪") || contains(weatherType, "大雾"):
		level = "高风险"
		tips = []string{
			"极端天气，建议暂缓出行",
			"如必须出行，请选择正规网约车平台",
			"行驶途中如遇危险请及时靠边停车",
			"紧急情况请拨打122交通报警电话",
		}
	case contains(weatherType, "雨") || contains(weatherType, "雪"):
		level = "中风险"
		tips = []string{
			"路面湿滑，请系好安全带",
			"建议司机打开车灯，保持安全车距",
			"避免急刹车和急转弯",
			"预留充足出行时间",
		}
	case contains(weatherType, "雾") || contains(weatherType, "霾"):
		level = "中风险"
		tips = []string{
			"能见度较低，请注意行车安全",
			"建议佩戴口罩保护呼吸道",
			"司机请开启雾灯低速行驶",
			"如能见度过低建议延迟出行",
		}
	case contains(weatherType, "大风") || weather.Windpower >= "6":
		level = "中风险"
		tips = []string{
			"大风天气，请注意高空坠物",
			"远离广告牌、大树等",
			"车辆行驶时注意横风影响",
			"建议关闭车窗行驶",
		}
	default:
		level = "低风险"
		tips = []string{
			"天气良好，适合出行",
			"请系好安全带，注意交通安全",
			"文明乘车，相互尊重",
			"祝您旅途愉快",
		}
	}

	return tips, level
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// -------------------------- 路线天气预警服务 --------------------------

// RouteWeatherWarning 路线天气预警结果
type RouteWeatherWarning struct {
	StartCity      string   `json:"start_city"`      // 起点城市
	EndCity        string   `json:"end_city"`        // 终点城市
	StartWeather   string   `json:"start_weather"`   // 起点天气
	EndWeather     string   `json:"end_weather"`     // 终点天气
	WarningLevel   string   `json:"warning_level"`   // 预警等级
	WarningMessage string   `json:"warning_message"` // 预警消息
	RouteAdvice    string   `json:"route_advice"`    // 路线建议
	SafetyTips     []string `json:"safety_tips"`     // 安全提示
}

// GetRouteWeatherWarning 获取路线天气预警
func GetRouteWeatherWarning(startAddress, endAddress string) (*RouteWeatherWarning, error) {
	// 获取起点天气
	startGeo, err := Geocode(startAddress)
	if err != nil {
		return nil, fmt.Errorf("起点地址解析失败: %v", err)
	}
	startWeather, err := WeatherQuery(startGeo.Adcode)
	if err != nil {
		return nil, fmt.Errorf("起点天气查询失败: %v", err)
	}

	// 获取终点天气
	endGeo, err := Geocode(endAddress)
	if err != nil {
		return nil, fmt.Errorf("终点地址解析失败: %v", err)
	}
	endWeather, err := WeatherQuery(endGeo.Adcode)
	if err != nil {
		return nil, fmt.Errorf("终点天气查询失败: %v", err)
	}

	// 分析路线天气预警
	warningLevel, warningMsg := analyzeRouteWeather(startWeather, endWeather)
	routeAdvice := generateRouteAdvice(startWeather, endWeather)
	safetyTips, _ := generateCombinedSafetyInfo(startWeather, endWeather)

	return &RouteWeatherWarning{
		StartCity:      startWeather.City,
		EndCity:        endWeather.City,
		StartWeather:   startWeather.Weather,
		EndWeather:     endWeather.Weather,
		WarningLevel:   warningLevel,
		WarningMessage: warningMsg,
		RouteAdvice:    routeAdvice,
		SafetyTips:     safetyTips,
	}, nil
}

// analyzeRouteWeather 分析路线天气预警等级
func analyzeRouteWeather(startWeather, endWeather *WeatherResult) (string, string) {
	dangerousWeathers := []string{"暴雨", "暴雪", "大雾", "台风", "冰雹"}
	moderateWeathers := []string{"雨", "雪", "雾", "霾"}

	// 检查是否有极端天气
	for _, dw := range dangerousWeathers {
		if contains(startWeather.Weather, dw) || contains(endWeather.Weather, dw) {
			return "红色预警", fmt.Sprintf("路线途经区域有%s天气，强烈建议延迟出行或改变路线", dw)
		}
	}

	// 检查是否有中等风险天气
	for _, mw := range moderateWeathers {
		if contains(startWeather.Weather, mw) || contains(endWeather.Weather, mw) {
			return "黄色预警", fmt.Sprintf("路线途经区域有%s天气，请注意行车安全，建议减速慢行", mw)
		}
	}

	return "绿色", "路线天气良好，适合出行"
}

// generateRouteAdvice 生成路线建议
func generateRouteAdvice(startWeather, endWeather *WeatherResult) string {
	if startWeather.Weather == endWeather.Weather {
		return fmt.Sprintf("起点和终点天气相同（%s），全程天气状况一致。", startWeather.Weather)
	}
	return fmt.Sprintf("起点%s（%s），终点%s（%s）。请注意天气变化，做好应对准备。",
		startWeather.City, startWeather.Weather, endWeather.City, endWeather.Weather)
}

// generateCombinedSafetyInfo 生成综合安全信息
func generateCombinedSafetyInfo(startWeather, endWeather *WeatherResult) ([]string, string) {
	tips1, level1 := generateSafetyInfo(startWeather)
	tips2, level2 := generateSafetyInfo(endWeather)

	// 合并去重安全提示
	tipMap := make(map[string]bool)
	var combinedTips []string
	for _, t := range tips1 {
		if !tipMap[t] {
			tipMap[t] = true
			combinedTips = append(combinedTips, t)
		}
	}
	for _, t := range tips2 {
		if !tipMap[t] {
			tipMap[t] = true
			combinedTips = append(combinedTips, t)
		}
	}

	// 取较高风险等级
	combinedLevel := level1
	if level2 == "高风险" || (level2 == "中风险" && level1 == "低风险") {
		combinedLevel = level2
	}

	return combinedTips, combinedLevel
}

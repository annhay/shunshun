// Package utils 包含地图相关的工具函数，封装了高德地图API的调用
package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"shunshun/internal/pkg/global"
)

// -------------------------- 高德地图API配置 --------------------------

// getAmapAPIKey 获取高德地图API密钥
// 从全局配置中获取，如果配置不存在则返回空字符串
func getAmapAPIKey() string {
	if global.AppConf != nil {
		return global.AppConf.Amap.APIKey
	}
	return "" // 默认值，防止nil指针异常
}

// getAmapBaseURL 获取高德地图API基础URL
// 从全局配置中获取，如果配置不存在则返回默认URL
func getAmapBaseURL() string {
	if global.AppConf != nil && global.AppConf.Amap.BaseURL != "" {
		return global.AppConf.Amap.BaseURL
	}
	return "https://restapi.amap.com" // 默认高德地图API URL
}

// -------------------------- API URL构建函数 --------------------------

// buildGeocodeURL 构建地理编码API URL
// 地理编码：将地址转换为经纬度坐标
func buildGeocodeURL() string {
	return getAmapBaseURL() + "/v3/geocode/geo"
}

// buildRegeocodeURL 构建逆地理编码API URL
// 逆地理编码：将经纬度转换为地址信息
func buildRegeocodeURL() string {
	return getAmapBaseURL() + "/v3/geocode/regeo"
}

// buildDrivingURL 构建驾车路径规划API URL
// 驾车路径规划：获取两点之间的驾车路线
func buildDrivingURL() string {
	return getAmapBaseURL() + "/v3/direction/driving"
}

// buildDistrictURL 构建行政区域查询API URL
// 行政区域查询：获取城市、区县等行政区域信息
func buildDistrictURL() string {
	return getAmapBaseURL() + "/v3/config/district"
}

// buildIPLocationURL 构建IP定位API URL
// IP定位：根据IP地址获取地理位置信息
func buildIPLocationURL() string {
	return getAmapBaseURL() + "/v3/ip"
}

// buildPOISearchURL 构建POI搜索API URL
// POI搜索：根据关键词搜索兴趣点（如餐厅、酒店等）
func buildPOISearchURL() string {
	return getAmapBaseURL() + "/v3/place/text"
}

// buildWeatherQueryURL 构建天气查询API URL
// 天气查询：获取指定城市的实时天气信息
func buildWeatherQueryURL() string {
	return getAmapBaseURL() + "/v3/weather/weatherInfo"
}

// -------------------------- 地理编码相关 --------------------------

// GeocodeResponse 地理编码响应结构
type GeocodeResponse struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Geocodes []struct {
		FormattedAddress string `json:"formatted_address"`
		Province         string `json:"province"`
		City             string `json:"city"`
		District         string `json:"district"`
		Adcode           string `json:"adcode"`
		Location         string `json:"location"` // 经纬度，格式：lng,lat
	} `json:"geocodes"`
}

// GeocodeResult 地理编码结果
type GeocodeResult struct {
	Address  string `json:"address"`
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
	Adcode   string `json:"adcode"`
	Lng      string `json:"lng"`
	Lat      string `json:"lat"`
}

// Geocode 地址转经纬度
func Geocode(address string) (*GeocodeResult, error) {
	encodedAddr := url.QueryEscape(address)
	reqURL := fmt.Sprintf("%s?address=%s&key=%s", buildGeocodeURL(), encodedAddr, getAmapAPIKey())

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res GeocodeResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if res.Status != "1" {
		return nil, fmt.Errorf("地理编码失败: %s", res.Info)
	}
	if len(res.Geocodes) == 0 {
		return nil, fmt.Errorf("未找到地理编码结果")
	}

	geo := res.Geocodes[0]
	// 解析经纬度，格式：lng,lat
	lng, lat := "", ""
	if geo.Location != "" {
		parts := strings.Split(geo.Location, ",")
		if len(parts) == 2 {
			lng = parts[0]
			lat = parts[1]
		}
	}

	return &GeocodeResult{
		Address:  geo.FormattedAddress,
		Province: geo.Province,
		City:     geo.City,
		District: geo.District,
		Adcode:   geo.Adcode,
		Lng:      lng,
		Lat:      lat,
	}, nil
}

// GeocodeWithCache 带缓存的地址转经纬度
func GeocodeWithCache(address string) (*GeocodeResult, error) {
	// 生成缓存键
	cacheKey := fmt.Sprintf("geocode:%s", url.QueryEscape(address))
	cacheExpiration := 24 * time.Hour

	// 检查缓存
	if global.Rdb != nil {
		cachedData, err := global.Rdb.Get(global.Ctx, cacheKey).Result()
		if err == nil {
			// 缓存命中
			var result GeocodeResult
			if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
				return &result, nil
			}
		}
	}

	// 缓存未命中，调用原始Geocode函数
	result, err := Geocode(address)
	if err != nil {
		return nil, err
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

// -------------------------- 逆地理编码相关 --------------------------

// RegeocodeResponse 逆地理编码响应结构
type RegeocodeResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	Regeocode struct {
		FormattedAddress string `json:"formatted_address"`
		AddressComponent struct {
			Province string `json:"province"`
			City     string `json:"city"`
			District string `json:"district"`
			Township string `json:"township"`
			Adcode   string `json:"adcode"`
		} `json:"addressComponent"`
	} `json:"regeocode"`
}

// RegeocodeResult 逆地理编码结果
type RegeocodeResult struct {
	Address  string `json:"address"`
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
	Township string `json:"township"`
	Adcode   string `json:"adcode"`
}

// Regeocode 经纬度转地址
func Regeocode(lng, lat string) (*RegeocodeResult, error) {
	reqURL := fmt.Sprintf("%s?location=%s,%s&key=%s", buildRegeocodeURL(), lng, lat, getAmapAPIKey())

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res RegeocodeResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if res.Status != "1" {
		return nil, fmt.Errorf("逆地理编码失败: %s", res.Info)
	}

	return &RegeocodeResult{
		Address:  res.Regeocode.FormattedAddress,
		Province: res.Regeocode.AddressComponent.Province,
		City:     res.Regeocode.AddressComponent.City,
		District: res.Regeocode.AddressComponent.District,
		Township: res.Regeocode.AddressComponent.Township,
		Adcode:   res.Regeocode.AddressComponent.Adcode,
	}, nil
}

// -------------------------- 驾车路径规划相关 --------------------------

// DrivingResponse 驾车路径规划响应结构
type DrivingResponse struct {
	Status string `json:"status"`
	Info   string `json:"info"`
	Route  struct {
		Origin      string `json:"origin"`
		Destination string `json:"destination"`
		Paths       []struct {
			Distance string `json:"distance"` // 距离（米）
			Duration string `json:"duration"` // 耗时（秒）
			Strategy string `json:"strategy"` // 策略
			Steps    []struct {
				Instruction string `json:"instruction"` // 导航指令
				Road        string `json:"road"`        // 道路名称
				Distance    string `json:"distance"`    // 距离
				Duration    string `json:"duration"`    // 耗时
			} `json:"steps"`
		} `json:"paths"`
	} `json:"route"`
}

// DrivingResult 驾车路径规划结果
type DrivingResult struct {
	Distance string              `json:"distance"` // 总距离（米）
	Duration string              `json:"duration"` // 总耗时（秒）
	Steps    []DrivingStepResult `json:"steps"`    // 路径步骤
}

// DrivingStepResult 路径步骤结果
type DrivingStepResult struct {
	Instruction string `json:"instruction"` // 导航指令
	Road        string `json:"road"`        // 道路名称
	Distance    string `json:"distance"`    // 距离
	Duration    string `json:"duration"`    // 耗时
}

// DrivingRoute 驾车路径规划（起点→终点）
func DrivingRoute(originLng, originLat, destLng, destLat string) (*DrivingResult, error) {
	origin := fmt.Sprintf("%s,%s", originLng, originLat)
	dest := fmt.Sprintf("%s,%s", destLng, destLat)
	reqURL := fmt.Sprintf("%s?origin=%s&destination=%s&key=%s", buildDrivingURL(), origin, dest, getAmapAPIKey())

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res DrivingResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if res.Status != "1" {
		return nil, fmt.Errorf("路径规划失败: %s", res.Info)
	}
	if len(res.Route.Paths) == 0 {
		return nil, fmt.Errorf("未找到路径")
	}

	path := res.Route.Paths[0]
	var steps []DrivingStepResult
	for _, s := range path.Steps {
		steps = append(steps, DrivingStepResult{
			Instruction: s.Instruction,
			Road:        s.Road,
			Distance:    s.Distance,
			Duration:    s.Duration,
		})
	}

	return &DrivingResult{
		Distance: path.Distance,
		Duration: path.Duration,
		Steps:    steps,
	}, nil
}

// -------------------------- 行政区域查询相关 --------------------------

// DistrictResponse 行政区域查询响应结构
type DistrictResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	Districts []struct {
		Name   string `json:"name"`
		Adcode string `json:"adcode"`
		Center string `json:"center"` // 区域中心点经纬度
		Level  string `json:"level"`  // 行政区级别
	} `json:"districts"`
}

// DistrictResult 行政区域查询结果
type DistrictResult struct {
	Name   string `json:"name"`
	Adcode string `json:"adcode"`
	Center string `json:"center"`
	Level  string `json:"level"`
}

// DistrictQuery 行政区域查询（支持按名称/编码查询）
func DistrictQuery(keyword string) ([]DistrictResult, error) {
	reqURL := fmt.Sprintf("%s?keywords=%s&key=%s", buildDistrictURL(), url.QueryEscape(keyword), getAmapAPIKey())
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res DistrictResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if res.Status != "1" {
		return nil, fmt.Errorf("行政区域查询失败: %s", res.Info)
	}

	var districts []DistrictResult
	for _, d := range res.Districts {
		districts = append(districts, DistrictResult{
			Name:   d.Name,
			Adcode: d.Adcode,
			Center: d.Center,
			Level:  d.Level,
		})
	}
	return districts, nil
}

// -------------------------- IP定位相关 --------------------------

// IPLocationResponse IP定位响应结构
type IPLocationResponse struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Province string `json:"province"`
	City     string `json:"city"`
	Adcode   string `json:"adcode"`
}

// IPLocationResult IP定位结果
type IPLocationResult struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Adcode   string `json:"adcode"`
}

// IPLocation IP定位（获取IP对应的城市）
func IPLocation(ip string) (*IPLocationResult, error) {
	reqURL := fmt.Sprintf("%s?ip=%s&key=%s", buildIPLocationURL(), ip, getAmapAPIKey())
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res IPLocationResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if res.Status != "1" {
		return nil, fmt.Errorf("IP定位失败: %s", res.Info)
	}

	return &IPLocationResult{
		Province: res.Province,
		City:     res.City,
		Adcode:   res.Adcode,
	}, nil
}

// -------------------------- POI搜索相关 --------------------------

// POISearchResponse POI搜索响应结构
type POISearchResponse struct {
	Status string `json:"status"`
	Info   string `json:"info"`
	Pois   []struct {
		Name     string `json:"name"`     // POI名称
		Location string `json:"location"` // 经纬度
		Address  string `json:"address"`  // 地址
		Type     string `json:"type"`     // 类型
		Tel      string `json:"tel"`      // 电话
		Distance string `json:"distance"` // 距离
	} `json:"pois"`
}

// POIResult POI搜索结果
type POIResult struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Address  string `json:"address"`
	Type     string `json:"type"`
	Tel      string `json:"tel"`
	Distance string `json:"distance"`
}

// POISearch POI搜索（关键词+城市）
func POISearch(keyword, city string) ([]POIResult, error) {
	reqURL := fmt.Sprintf("%s?keywords=%s&city=%s&key=%s", buildPOISearchURL(), url.QueryEscape(keyword), url.QueryEscape(city), getAmapAPIKey())
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res POISearchResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if res.Status != "1" {
		return nil, fmt.Errorf("POI搜索失败: %s", res.Info)
	}

	var pois []POIResult
	for _, p := range res.Pois {
		pois = append(pois, POIResult{
			Name:     p.Name,
			Location: p.Location,
			Address:  p.Address,
			Type:     p.Type,
			Tel:      p.Tel,
			Distance: p.Distance,
		})
	}
	return pois, nil
}

// -------------------------- 天气查询相关 --------------------------

// WeatherResponse 天气查询响应结构
type WeatherResponse struct {
	Status string `json:"status"`
	Info   string `json:"info"`
	Lives  []struct {
		Province         string `json:"province"`
		City             string `json:"city"`
		Adcode           string `json:"adcode"`
		Weather          string `json:"weather"`       // 天气现象（如"晴"）
		Temperature      string `json:"temperature"`   // 实时温度
		Winddirection    string `json:"winddirection"` // 风向
		Windpower        string `json:"windpower"`     // 风力
		Humidity         string `json:"humidity"`      // 湿度
		Reporttime       string `json:"reporttime"`    // 数据发布时间
		TemperatureFloat string `json:"temperature_float"`
		HumidityFloat    string `json:"humidity_float"`
	} `json:"lives"`
	Forecasts []struct {
		Province   string `json:"province"`
		City       string `json:"city"`
		Adcode     string `json:"adcode"`
		Reporttime string `json:"reporttime"`
		Casts      []struct {
			Date         string `json:"date"`
			Week         string `json:"week"`
			Dayweather   string `json:"dayweather"`   // 白天天气
			Nightweather string `json:"nightweather"` // 夜间天气
			Daytemp      string `json:"daytemp"`      // 白天温度
			Nighttemp    string `json:"nighttemp"`    // 夜间温度
			Daywind      string `json:"daywind"`      // 白天风向
			Nightwind    string `json:"nightwind"`    // 夜间风向
			Daypower     string `json:"daypower"`     // 白天风力
			Nightpower   string `json:"nightpower"`   // 夜间风力
		} `json:"casts"`
	} `json:"forecasts"`
}

// WeatherResult 实时天气结果
type WeatherResult struct {
	Province      string `json:"province"`
	City          string `json:"city"`
	Adcode        string `json:"adcode"`
	Weather       string `json:"weather"`       // 天气现象
	Temperature   string `json:"temperature"`   // 温度
	Winddirection string `json:"winddirection"` // 风向
	Windpower     string `json:"windpower"`     // 风力
	Humidity      string `json:"humidity"`      // 湿度
	Reporttime    string `json:"reporttime"`    // 发布时间
}

// WeatherForecastResult 天气预报结果
type WeatherForecastResult struct {
	Province   string              `json:"province"`
	City       string              `json:"city"`
	Reporttime string              `json:"reporttime"`
	Casts      []WeatherCastResult `json:"casts"`
}

// WeatherCastResult 单日天气预报
type WeatherCastResult struct {
	Date         string `json:"date"`
	Week         string `json:"week"`
	Dayweather   string `json:"dayweather"`
	Nightweather string `json:"nightweather"`
	Daytemp      string `json:"daytemp"`
	Nighttemp    string `json:"nighttemp"`
	Daywind      string `json:"daywind"`
	Nightwind    string `json:"nightwind"`
	Daypower     string `json:"daypower"`
	Nightpower   string `json:"nightpower"`
}

// WeatherQuery 实时天气查询（按城市编码）
func WeatherQuery(adcode string) (*WeatherResult, error) {
	reqURL := fmt.Sprintf("%s?city=%s&key=%s&extensions=base", buildWeatherQueryURL(), adcode, getAmapAPIKey())
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res WeatherResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if res.Status != "1" {
		return nil, fmt.Errorf("天气查询失败: %s", res.Info)
	}
	if len(res.Lives) == 0 {
		return nil, fmt.Errorf("未获取到天气数据")
	}

	life := res.Lives[0]
	return &WeatherResult{
		Province:      life.Province,
		City:          life.City,
		Adcode:        life.Adcode,
		Weather:       life.Weather,
		Temperature:   life.Temperature,
		Winddirection: life.Winddirection,
		Windpower:     life.Windpower,
		Humidity:      life.Humidity,
		Reporttime:    life.Reporttime,
	}, nil
}

// WeatherForecast 天气预报查询（按城市编码，返回未来几天预报）
func WeatherForecast(adcode string) (*WeatherForecastResult, error) {
	reqURL := fmt.Sprintf("%s?city=%s&key=%s&extensions=all", buildWeatherQueryURL(), adcode, getAmapAPIKey())
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res WeatherResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if res.Status != "1" {
		return nil, fmt.Errorf("天气预报查询失败: %s", res.Info)
	}
	if len(res.Forecasts) == 0 {
		return nil, fmt.Errorf("未获取到天气预报数据")
	}

	forecast := res.Forecasts[0]
	var casts []WeatherCastResult
	for _, c := range forecast.Casts {
		casts = append(casts, WeatherCastResult{
			Date:         c.Date,
			Week:         c.Week,
			Dayweather:   c.Dayweather,
			Nightweather: c.Nightweather,
			Daytemp:      c.Daytemp,
			Nighttemp:    c.Nighttemp,
			Daywind:      c.Daywind,
			Nightwind:    c.Nightwind,
			Daypower:     c.Daypower,
			Nightpower:   c.Nightpower,
		})
	}

	return &WeatherForecastResult{
		Province:   forecast.Province,
		City:       forecast.City,
		Reporttime: forecast.Reporttime,
		Casts:      casts,
	}, nil
}

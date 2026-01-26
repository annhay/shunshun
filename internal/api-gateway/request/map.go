package request

// -------------------------- 地理编码相关请求 --------------------------

// GeocodeReq 地理编码请求
type GeocodeReq struct {
	Address string `form:"address" json:"address" binding:"required"` // 地址
}

// RegeocodeReq 逆地理编码请求
type RegeocodeReq struct {
	Lng string `form:"lng" json:"lng" binding:"required"` // 经度
	Lat string `form:"lat" json:"lat" binding:"required"` // 纬度
}

// -------------------------- 路径规划相关请求 --------------------------

// DrivingRouteReq 驾车路径规划请求
type DrivingRouteReq struct {
	OriginLng string `form:"origin_lng" json:"origin_lng" binding:"required"` // 起点经度
	OriginLat string `form:"origin_lat" json:"origin_lat" binding:"required"` // 起点纬度
	DestLng   string `form:"dest_lng" json:"dest_lng" binding:"required"`     // 终点经度
	DestLat   string `form:"dest_lat" json:"dest_lat" binding:"required"`     // 终点纬度
}

// DrivingRouteByAddressReq 按地址驾车路径规划请求
type DrivingRouteByAddressReq struct {
	Origin      string `form:"origin" json:"origin" binding:"required"`           // 起点地址
	Destination string `form:"destination" json:"destination" binding:"required"` // 终点地址
}

// -------------------------- 行政区域查询请求 --------------------------

// DistrictQueryReq 行政区域查询请求
type DistrictQueryReq struct {
	Keyword string `form:"keyword" json:"keyword" binding:"required"` // 关键词（城市名/区县名/编码）
}

// -------------------------- IP定位请求 --------------------------

// IPLocationReq IP定位请求
type IPLocationReq struct {
	IP string `form:"ip" json:"ip"` // IP地址（可选，不传则使用请求IP）
}

// -------------------------- POI搜索请求 --------------------------

// POISearchReq POI搜索请求
type POISearchReq struct {
	Keyword string `form:"keyword" json:"keyword" binding:"required"` // 搜索关键词
	City    string `form:"city" json:"city" binding:"required"`       // 城市名称或编码
}

// -------------------------- 天气查询请求 --------------------------

// WeatherQueryReq 天气查询请求
type WeatherQueryReq struct {
	Adcode string `form:"adcode" json:"adcode" binding:"required"` // 城市行政编码
}

// WeatherForecastReq 天气预报请求
type WeatherForecastReq struct {
	Adcode string `form:"adcode" json:"adcode" binding:"required"` // 城市行政编码
}

// -------------------------- AI天气提醒请求 --------------------------

// WeatherReminderReq AI天气提醒请求（按城市编码）
type WeatherReminderReq struct {
	Adcode string `form:"adcode" json:"adcode" binding:"required"` // 城市行政编码
}

// WeatherReminderByAddressReq AI天气提醒请求（按地址）
type WeatherReminderByAddressReq struct {
	Address string `form:"address" json:"address" binding:"required"` // 地址
}

// WeatherReminderByLocationReq AI天气提醒请求（按经纬度）
type WeatherReminderByLocationReq struct {
	Lng string `form:"lng" json:"lng" binding:"required"` // 经度
	Lat string `form:"lat" json:"lat" binding:"required"` // 纬度
}

// -------------------------- 路线天气预警请求 --------------------------

// RouteWeatherWarningReq 路线天气预警请求
type RouteWeatherWarningReq struct {
	StartAddress string `form:"start_address" json:"start_address" binding:"required"` // 起点地址
	EndAddress   string `form:"end_address" json:"end_address" binding:"required"`     // 终点地址
}

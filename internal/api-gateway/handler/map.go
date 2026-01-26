package handler

import (
	"net/http"
	"shunshun/internal/api-gateway/request"
	"shunshun/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// -------------------------- 地理编码相关Handler --------------------------

// Geocode 地理编码（地址转经纬度）
func Geocode(c *gin.Context) {
	var form request.GeocodeReq
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	result, err := utils.Geocode(form.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": result})
}

// Regeocode 逆地理编码（经纬度转地址）
func Regeocode(c *gin.Context) {
	var form request.RegeocodeReq
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	result, err := utils.Regeocode(form.Lng, form.Lat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": result})
}

// -------------------------- 路径规划相关Handler --------------------------

// DrivingRoute 驾车路径规划
func DrivingRoute(c *gin.Context) {
	var form request.DrivingRouteReq
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	result, err := utils.DrivingRoute(form.OriginLng, form.OriginLat, form.DestLng, form.DestLat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": result})
}

// DrivingRouteByAddress 按地址驾车路径规划
func DrivingRouteByAddress(c *gin.Context) {
	var form request.DrivingRouteByAddressReq
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 解析起点地址
	originGeo, err := utils.Geocode(form.Origin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "起点地址解析失败: " + err.Error()})
		return
	}

	// 解析终点地址
	destGeo, err := utils.Geocode(form.Destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "终点地址解析失败: " + err.Error()})
		return
	}

	// 规划路线
	result, err := utils.DrivingRoute(originGeo.Lng, originGeo.Lat, destGeo.Lng, destGeo.Lat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"origin":      originGeo,
			"destination": destGeo,
			"route":       result,
		},
	})
}

// -------------------------- 行政区域查询Handler --------------------------

// DistrictQuery 行政区域查询
func DistrictQuery(c *gin.Context) {
	var form request.DistrictQueryReq
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	result, err := utils.DistrictQuery(form.Keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": result})
}

// -------------------------- IP定位Handler --------------------------

// IPLocation IP定位
func IPLocation(c *gin.Context) {
	var form request.IPLocationReq
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 如果未传IP，使用请求客户端IP
	ip := form.IP
	if ip == "" {
		ip = c.ClientIP()
	}

	result, err := utils.IPLocation(ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": result})
}

// -------------------------- POI搜索Handler --------------------------

// POISearch POI搜索
func POISearch(c *gin.Context) {
	var form request.POISearchReq
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	result, err := utils.POISearch(form.Keyword, form.City)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": result})
}

// -------------------------- 天气查询Handler --------------------------

// WeatherQuery 实时天气查询
func WeatherQuery(c *gin.Context) {
	var form request.WeatherQueryReq
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	result, err := utils.WeatherQuery(form.Adcode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": result})
}

// WeatherForecast 天气预报查询
func WeatherForecast(c *gin.Context) {
	var form request.WeatherForecastReq
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	result, err := utils.WeatherForecast(form.Adcode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": result})
}

// -------------------------- AI天气提醒Handler --------------------------

// WeatherReminder AI智能天气提醒（按城市编码）
func WeatherReminder(c *gin.Context) {
	var form request.WeatherReminderReq
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	result, err := utils.GetWeatherReminder(form.Adcode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": result})
}

// WeatherReminderByAddress AI智能天气提醒（按地址）
func WeatherReminderByAddress(c *gin.Context) {
	var form request.WeatherReminderByAddressReq
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	result, err := utils.GetWeatherReminderByAddress(form.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": result})
}

// WeatherReminderByLocation AI智能天气提醒（按经纬度）
func WeatherReminderByLocation(c *gin.Context) {
	var form request.WeatherReminderByLocationReq
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	result, err := utils.GetWeatherReminderByLocation(form.Lng, form.Lat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": result})
}

// RouteWeatherWarning 路线天气预警
func RouteWeatherWarning(c *gin.Context) {
	var form request.RouteWeatherWarningReq
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	result, err := utils.GetRouteWeatherWarning(form.StartAddress, form.EndAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": result})
}

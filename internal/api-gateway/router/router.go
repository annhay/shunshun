package router

import (
	"shunshun/internal/api-gateway/consts"
	"shunshun/internal/api-gateway/handler"
	"shunshun/internal/api-gateway/middleware"

	"github.com/gin-gonic/gin"
)

// LoadRouter 加载路由配置
//
// 返回值:
//   - *gin.Engine: Gin引擎实例
//
// 路由分组:
//  1. 公共接口(/api/v1/public): 无需身份验证
//  2. 私有接口(/api/v1/private): 需要JWT身份验证
func LoadRouter() *gin.Engine {
	router := gin.Default()

	// 使用跨域中间件
	router.Use(middleware.Cors())

	// 公共接口,无需中间件验证
	public := router.Group("/api/v1/public")
	{
		// 用户相关接口
		public.POST("/sendTextMessage", handler.SendTextMessage) // 短信发送
		public.POST("/register", handler.Register)               // 注册
		public.POST("/login", handler.Login)                     // 登录
		public.POST("/forgotPassword", handler.ForgotPassword)   // 修改密码

		// 地图服务接口
		mapGroup := public.Group("/map")
		{
			mapGroup.POST("/geocode", handler.Geocode)                       // 地理编码
			mapGroup.POST("/regeocode", handler.Regeocode)                   // 逆地理编码
			mapGroup.POST("/driving", handler.DrivingRoute)                  // 驾车路径规划
			mapGroup.POST("/driving/address", handler.DrivingRouteByAddress) // 按地址驾车路径规划
			mapGroup.POST("/district", handler.DistrictQuery)                // 行政区域查询
			mapGroup.POST("/ip", handler.IPLocation)                         // IP 定位
			mapGroup.POST("/poi", handler.POISearch)                         // POI 搜索
		}

		// 天气服务接口
		weatherGroup := public.Group("/weather")
		{
			weatherGroup.POST("/now", handler.WeatherQuery)         // 实时天气查询
			weatherGroup.POST("/forecast", handler.WeatherForecast) // 天气预报查询

			// AI智能天气提醒
			aiGroup := weatherGroup.Group("/ai")
			{
				aiGroup.POST("/reminder", handler.WeatherReminder)                    // AI天气提醒(按城市编码)
				aiGroup.POST("/reminder/address", handler.WeatherReminderByAddress)   // AI天气提醒(按地址)
				aiGroup.POST("/reminder/location", handler.WeatherReminderByLocation) // AI天气提醒(按经纬度)
				aiGroup.POST("/route", handler.RouteWeatherWarning)                   // 路线天气预警
			}
		}
	}

	// 私有接口,需用户鉴权
	private := router.Group("/api/v1/private")
	private.Use(middleware.ParseToken(consts.JwtKey)) // 使用JWT验证中间件
	{
		// 用户相关接口
		private.POST("/completeInformation", handler.CompleteInformation) // 完善信息
		private.POST("/studentVerification", handler.StudentVerification) // 学生认证
		private.GET("/personalCenter", handler.PersonalCenter)            // 个人中心
		private.POST("/logout", handler.Logout)                           // 账号注销

		// 司机相关接口
		private.POST("/newDriver", handler.NewDriver) // 司机信息添加
		private.POST("/updDriver", handler.UpdDriver) // 司机信息修改

		// 订单相关接口
		private.POST("/newOrder", handler.NewOrder)         // 创建订单
		private.POST("/acceptOrders", handler.AcceptOrders) // 司机接单
	}

	return router
}

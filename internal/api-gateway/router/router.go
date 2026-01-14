package router

import (
	"shunshun/internal/api-gateway/consts"
	"shunshun/internal/api-gateway/handler"
	"shunshun/internal/api-gateway/middleware"

	"github.com/gin-gonic/gin"
)

func LoadRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.Cors())
	//公共接口,无需中间件验证
	public := router.Group("/api/v1/public")
	{
		public.POST("/sendTextMessage", handler.SendTextMessage) //短信发送
		public.POST("/register", handler.Register)               //注册
		public.POST("/login", handler.Login)                     //登录
		public.POST("/forgotPassword", handler.ForgotPassword)   //修改密码
	}
	//私有接口,需用户鉴权
	private := router.Group("/api/v1/private")
	private.Use(middleware.ParseToken(consts.JwtKey))
	{
		private.POST("/completeInformation", handler.CompleteInformation) //完善信息
		private.POST("/studentVerification", handler.StudentVerification) //学生认证
		private.GET("/personalCenter", handler.PersonalCenter)            //个人中心
		private.POST("/logout", handler.Logout)                           //账号注销
		private.POST("/newDriver", handler.NewDriver)                     //司机信息添加
		private.POST("/updDriver", handler.UpdDriver)                     //司机信息修改
	}
	return router
}

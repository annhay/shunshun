package router

import (
	"shunshun/internal/api-gateway/handler"
	"shunshun/internal/api-gateway/middleware"

	"github.com/gin-gonic/gin"
)

func LoadRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.Cors())
	//公共接口,无需中间件验证
	public := router.Group("/api/v1")
	{
		public.POST("/sendTextMessage", handler.SendTextMessage) //短信发送
	}
	//私有接口,需中间件验证
	private := router.Group("/api/v1")
	{
		private.POST("")
	}
	return router
}

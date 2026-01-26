package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
// 处理跨域资源共享(CORS)请求，允许不同域名的前端应用访问API
// 
// 返回值:
//   - gin.HandlerFunc: Gin框架的中间件处理函数
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		
		// 如果请求包含Origin头，设置CORS响应头
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")                                  // 允许所有来源
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")     // 允许的HTTP方法
			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization") // 允许的请求头
			c.Header("Access-Control-Allow-Credentials", "true")                           // 允许携带凭证
			c.Set("content-type", "application/json")                                      // 设置响应内容类型
		}
		
		// 处理预检请求(OPTIONS)
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // 返回204 No Content
		}
		
		// 继续处理请求
		c.Next()
	}
}

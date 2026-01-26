package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// MyCustomClaims JWT自定义声明结构体
// 包含用户ID和标准JWT声明

type MyCustomClaims struct {
	UserId uint64 `json:"userId"` // 用户ID
	jwt.RegisteredClaims          // 标准JWT声明
}

// CreateToken 创建JWT令牌
// 
// 参数:
//   - userId uint64: 用户ID
//   - key string: 签名密钥
// 
// 返回值:
//   - string: 生成的JWT令牌字符串，如果出错则返回空字符串
func CreateToken(userId uint64, key string) string {
	// 创建包含用户ID和过期时间的声明
	claims := MyCustomClaims{
		userId,
		jwt.RegisteredClaims{
			// 设置过期时间为30天
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		},
	}
	
	// 使用HS256算法创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 使用密钥签名令牌
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return ""
	}
	
	return ss
}

// ParseToken 解析JWT令牌的中间件
// 
// 参数:
//   - key string: 签名密钥
// 
// 返回值:
//   - gin.HandlerFunc: Gin框架的中间件处理函数
func ParseToken(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		tokenStr := c.Request.Header.Get("x-token")
		if tokenStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "账号未登录"})
			c.Abort()
			return
		}
		
		// 解析令牌
		token, err := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(key), nil
		})
		
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "token解析失败,无效token"})
			c.Abort()
			return
		}
		
		// 验证令牌声明
		claims, ok := token.Claims.(*MyCustomClaims)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "登录过期"})
			c.Abort()
			return
		}
		
		// 将用户ID存储到上下文
		c.Set("userId", claims.UserId)
	}
}

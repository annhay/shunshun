package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	UserId uint64 `json:"userId"`
	jwt.RegisteredClaims
}

// CreateToken 创建中间鉴权
func CreateToken(userId uint64, key string) string {
	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		userId,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return ""
	}
	return ss
}

// ParseToken 解析中间鉴权
func ParseToken(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("x-token")
		if tokenStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "账号未登录"})
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(key), nil
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "token解析失败,无效token"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(*MyCustomClaims)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "登录过期"})
			c.Abort()
			return
		}
		c.Set("userId", claims.UserId)
	}
}

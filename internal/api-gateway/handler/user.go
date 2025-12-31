package handler

import (
	"net/http"
	"shunshun/internal/api-gateway/request"
	"shunshun/internal/pkg/global"
	"shunshun/internal/proto"

	"github.com/gin-gonic/gin"
)

// SendTextMessage 短信发送
func SendTextMessage(c *gin.Context) {
	var form request.SendTextMessage
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := global.UserClient.SendTextMessage(c, &proto.SendTextMessageReq{Phone: form.Phone})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "短信发送成功"})
}

package handler

import (
	"net/http"
	"shunshun/internal/api-gateway/request"
	"shunshun/internal/pkg/global"
	"shunshun/internal/pkg/utils"
	"shunshun/internal/proto"

	"github.com/gin-gonic/gin"
)

// NewOrder 创建订单
//
// 参数:
//   - c *gin.Context: Gin上下文
//
// 处理逻辑:
//  1. 绑定请求参数
//  2. 调用订单服务创建订单
//  3. 返回创建结果
func NewOrder(c *gin.Context) {
	var form request.NewOrder
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := utils.Validate(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	resp, err := global.OrderClient.NewOrder(c, &proto.NewOrderReq{
		UserId:             form.UserId,
		TripType:           form.TripType,
		RideMode:           form.RideMode,
		PassengerNum:       int64(form.PassengerNum),
		CarType:            form.CarType,
		DepartureTime:      form.DepartureTime,
		WaitingTime:        int64(form.WaitingTime),
		StratDetailAddress: form.StratDetailAddress,
		EndDetailAddress:   form.EndDetailAddress,
		PaymentMethod:      form.PaymentMethod,
		CouponId:           form.CouponId,
		Remark:             form.Remark,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orderCode": resp.OrderCode,
		"payUrl":    resp.PayUrl,
	})
}

// AcceptOrders 司机接单
//
// 参数:
//   - c *gin.Context: Gin上下文
//
// 处理逻辑:
//  1. 绑定请求参数
//  2. 调用订单服务接单
//  3. 返回接单结果
func AcceptOrders(c *gin.Context) {
	var form request.AcceptOrders
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := utils.Validate(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	_, err := global.OrderClient.AcceptOrders(c, &proto.AcceptOrdersReq{
		UserId:  form.UserId,
		OrderId: form.OrderId,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "接单成功",
	})
}

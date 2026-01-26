package task

import (
	"shunshun/internal/pkg/global"
	"shunshun/internal/pkg/model"
	"time"

	"go.uber.org/zap"
)

// CheckUnpaidOrders 检查未支付的顺风车订单，10分钟未支付自动取消
func CheckUnpaidOrders() {
	global.Logger.Info("开始检查未支付的顺风车订单")

	// 查找10分钟前创建的未支付顺风车订单
	var orders []model.ShunOrder
	tenMinutesAgo := time.Now().Add(-10 * time.Minute)

	result := global.DB.Where("trip_type = ? AND payment_status = ? AND order_status = ? AND created_at < ?",
		"1", "1", "1", tenMinutesAgo).Find(&orders)

	if result.Error != nil {
		global.Logger.Error("查询未支付订单失败", zap.Error(result.Error))
		return
	}

	global.Logger.Info("查询到未支付顺风车订单", zap.Int("count", len(orders)))

	// 取消这些订单
	for _, order := range orders {
		global.Logger.Info("取消超时未支付的顺风车订单", zap.String("orderCode", order.OrderCode), zap.Int64("userId", order.UserId))

		// 更新订单状态为已取消
		if err := global.DB.Model(&order).Updates(map[string]interface{}{
			"order_status":   "4", // 已取消
			"payment_status": "4", // 支付失败
		}).Error; err != nil {
			global.Logger.Error("取消订单失败", zap.String("orderCode", order.OrderCode), zap.Error(err))
			continue
		}

		// 恢复优惠券（如果使用了优惠券）
		if order.CouponId > 0 {
			if err := global.DB.Model(&model.ShunCoupon{}).Where("id = ?", order.CouponId).Update("status", "1").Error; err != nil {
				global.Logger.Error("恢复优惠券失败", zap.Int64("couponId", order.CouponId), zap.Error(err))
			}
		}

		global.Logger.Info("订单取消成功", zap.String("orderCode", order.OrderCode))
	}
}

// StartOrderTaskScheduler 启动订单任务调度器
func StartOrderTaskScheduler() {
	global.Logger.Info("启动订单任务调度器")

	// 每10秒检查一次
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			CheckUnpaidOrders()
		}
	}()
}

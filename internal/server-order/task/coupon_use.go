package task

import (
	"errors"
	"shunshun/internal/pkg/global"
	"shunshun/internal/pkg/model"
	"time"
)

// CalculateCouponAmount 计算优惠券金额
func CalculateCouponAmount(userId int64, couponId int64, estimatedAmount float64) (float64, error) {
	if couponId <= 0 {
		return 0, nil
	}
	
	// 查询用户优惠券使用记录
	var useCoupon model.ShunUseCoupon
	result := global.DB.Where("user_id = ? AND coupon_id = ? AND use_status = 1", userId, couponId).First(&useCoupon)
	if result.Error != nil {
		return 0, errors.New("优惠券不存在或已使用")
	}
	
	// 检查优惠券是否过期
	if time.Now().After(useCoupon.ExpirationTime) {
		return 0, errors.New("优惠券已过期")
	}
	
	// 查询优惠券详情
	var coupon model.ShunCoupon
	if err := global.DB.First(&coupon, couponId).Error; err != nil {
		return 0, errors.New("优惠券不存在")
	}
	
	// 计算优惠金额
	var couponAmount float64
	switch coupon.CouponType {
	case "1": // 满减券
		if estimatedAmount >= coupon.SpendMoney {
			couponAmount = coupon.Money
		}
	case "2": // 折扣券
		discountAmount := estimatedAmount * float64(coupon.Discount) / 100
		if discountAmount > coupon.TotalMoney {
			couponAmount = coupon.TotalMoney
		} else {
			couponAmount = discountAmount
		}
	case "3": // 无门槛券
		couponAmount = coupon.Money
	}
	
	// 确保优惠金额不超过订单金额
	if couponAmount > estimatedAmount {
		couponAmount = estimatedAmount
	}
	
	return couponAmount, nil
}

// UpdateCouponStatus 更新优惠券状态
func UpdateCouponStatus(userId int64, couponId int64) error {
	if couponId <= 0 {
		return nil
	}
	
	// 更新优惠券状态
	if err := global.DB.Model(&model.ShunUseCoupon{}).Where("user_id = ? AND coupon_id = ?", userId, couponId).Updates(map[string]interface{}{
		"use_status": 2, // 已使用
		"use_time":   time.Now(),
	}).Error; err != nil {
		return err
	}
	
	return nil
}

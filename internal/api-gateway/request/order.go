package request

// NewOrder 订单创建请求
// @Summary 订单创建
// @Description 创建新订单
// @Tags 订单
// @Accept json
// @Produce json
// @Param request body NewOrder true "订单创建请求"
// @Success 200 {object} Response
// @Router /api/order/new [post]
type NewOrder struct {
	UserId             int64  `form:"user_id" json:"user_id" xml:"user_id" binding:"required,gt=0"`
	TripType           string `form:"trip_type" json:"trip_type" xml:"trip_type" binding:"required,oneof=1 2"`
	RideMode           string `form:"ride_mode" json:"ride_mode" xml:"ride_mode" binding:"required,oneof=1 2"`
	PassengerNum       int    `form:"passenger_num" json:"passenger_num" xml:"passenger_num" binding:"required,min=1,max=6"`
	CarType            string `form:"car_type" json:"car_type" xml:"car_type" binding:"required,oneof=1 2 3"`
	DepartureTime      string `form:"departure_time" json:"departure_time" xml:"departure_time" binding:"required"`
	WaitingTime        int    `form:"waiting_time" json:"waiting_time" xml:"waiting_time" binding:"required,min=10,max=180"`
	StratDetailAddress string `form:"strat_detail_address" json:"strat_detail_address" xml:"strat_detail_address" binding:"required"`
	EndDetailAddress   string `form:"end_detail_address" json:"end_detail_address" xml:"end_detail_address" binding:"required"`
	PaymentMethod      string `form:"payment_method" json:"payment_method" xml:"payment_method" binding:"required,oneof=1 2 3"`
	CouponId           int64  `form:"coupon_id" json:"coupon_id" xml:"coupon_id" binding:"omitempty,gt=0"`
	Remark             string `form:"remark" json:"remark" xml:"remark"`
}

// AcceptOrders 司机接单请求
// @Summary 司机接单
// @Description 司机接受订单
// @Tags 订单
// @Accept json
// @Produce json
// @Param request body AcceptOrders true "司机接单请求"
// @Success 200 {object} Response
// @Router /api/order/accept [post]
type AcceptOrders struct {
	UserId  int64 `form:"user_id" json:"user_id" xml:"user_id" binding:"required,gt=0"`
	OrderId int64 `form:"order_id" json:"order_id" xml:"order_id" binding:"required,gt=0"`
}

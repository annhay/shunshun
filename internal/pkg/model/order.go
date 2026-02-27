package model

import (
	"time"

	"gorm.io/gorm"
)

type ShunOrder struct { //订单信息表
	Id                    uint64    `gorm:"column:id;type:bigint UNSIGNED;primaryKey;not null;" json:"id"`
	OrderCode             string    `gorm:"column:order_code;type:varchar(64);comment:订单编号;not null;" json:"order_code"`                                                                           // 订单编号
	UserId                int64     `gorm:"column:user_id;type:bigint;comment:用户ID;not null;" json:"user_id"`                                                                                      // 用户ID
	TripType              string    `gorm:"column:trip_type;type:varchar(10);comment:行程模式：1-顺风车，2-打车，3-送货，4-宠物;not null;" json:"trip_type"`                                                        // 行程模式：1-顺风车，2-打车，3-送货，4-宠物
	RideMode              string    `gorm:"column:ride_mode;type:varchar(10);comment:乘车模式：1-拼座（拼单），2-只拼一单（拼单），3-独享【顺风车】，4-特惠拼送（送货），5-一对一专送（送货）【送货】，6-宠物快车（宠物）【宠物】;default:NULL;" json:"ride_mode"` // 乘车模式：1-拼座（拼单），2-只拼一单（拼单），3-独享【顺风车】，4-特惠拼送（送货），5-一对一专送（送货）【送货】，6-宠物快车（宠物）【宠物】
	PassengerNum          int32     `gorm:"column:passenger_num;type:int;comment:人数（最少1人，最多6人）【打车默认1人，不可设置】;default:1;" json:"passenger_num"`                                                      // 人数（最少1人，最多6人）【打车默认1人，不可设置】
	CarType               string    `gorm:"column:car_type;type:varchar(10);comment:车辆类型：1-经济车 2-商务车 3六座车;default:1;" json:"car_type"`                                                             // 车辆类型：1-经济车 2-商务车 3六座车
	DepartureTime         time.Time `gorm:"column:departure_time;type:datetime;comment:出发时间【顺风车】;default:NULL;" json:"departure_time"`                                                             // 出发时间【顺风车】
	WaitingTime           int32     `gorm:"column:waiting_time;type:int;comment:愿等待时间(分钟)（最短10分钟，最长3小时）【顺风车】;default:10;" json:"waiting_time"`                                                     // 愿等待时间(分钟)（最短10分钟，最长3小时）【顺风车】
	Remark                string    `gorm:"column:remark;type:varchar(255);comment:备注;default:NULL;" json:"remark"`                                                                                // 备注
	MainGroupOrderId      int64     `gorm:"column:main_group_order_id;type:bigint;comment:拼单主ID【拼单】;default:NULL;" json:"main_group_order_id"`                                                     // 拼单主ID【拼单】
	GroupUserCount        int32     `gorm:"column:group_user_count;type:int;comment:拼单人数【拼单，只在主订单显示】;default:0;" json:"group_user_count"`                                                          // 拼单人数【拼单，只在主订单显示】
	StartDetailAddress    string    `gorm:"column:start_detail_address;type:varchar(255);comment:起点详细地址;not null;" json:"start_detail_address"`                                                    // 起点详细地址
	StartLongitude        float64   `gorm:"column:start_longitude;type:decimal(10, 6);comment:起点经度;not null;" json:"start_longitude"`                                                              // 起点经度
	StartLatitude         float64   `gorm:"column:start_latitude;type:decimal(10, 6);comment:起点纬度;not null;" json:"start_latitude"`                                                                // 起点纬度
	EndDetailAddress      string    `gorm:"column:end_detail_address;type:varchar(255);comment:终点详细地址;not null;" json:"end_detail_address"`                                                        // 终点详细地址
	EndLongitude          float64   `gorm:"column:end_longitude;type:decimal(10, 6);comment:终点经度;not null;" json:"end_longitude"`                                                                  // 终点经度
	EndLatitude           float64   `gorm:"column:end_latitude;type:decimal(10, 6);comment:终点纬度;not null;" json:"end_latitude"`                                                                    // 终点纬度
	TravelDistance        float64   `gorm:"column:travel_distance;type:decimal(10, 2);comment:行程距离（公里）;default:NULL;" json:"travel_distance"`                                                      // 行程距离（公里）
	TravelDuration        int32     `gorm:"column:travel_duration;type:int;comment:行程时长;default:NULL;" json:"travel_duration"`                                                                     // 行程时长
	DriverId              int64     `gorm:"column:driver_id;type:bigint;comment:司机ID;default:NULL;" json:"driver_id"`                                                                              // 司机ID
	CarId                 int64     `gorm:"column:car_id;type:bigint;comment:车辆ID;default:NULL;" json:"car_id"`                                                                                    // 车辆ID
	ActualPassengerNum    int32     `gorm:"column:actual_passenger_num;type:int;comment:实际乘车人数（司机填写）【打车】;default:NULL;" json:"actual_passenger_num"`                                               // 实际乘车人数（司机填写）【打车】
	DriverAcceptAddress   string    `gorm:"column:driver_accept_address;type:varchar(255);comment:司机接单地址;default:NULL;" json:"driver_accept_address"`                                              // 司机接单地址
	DriverAcceptLongitude float64   `gorm:"column:driver_accept_longitude;type:decimal(10, 6);comment:司机接单经度;default:NULL;" json:"driver_accept_longitude"`                                        // 司机接单经度
	DriverAcceptLatitude  float64   `gorm:"column:driver_accept_latitude;type:decimal(10, 6);comment:司机接单纬度;default:NULL;" json:"driver_accept_latitude"`                                          // 司机接单纬度
	EstimatedAmount       float64   `gorm:"column:estimated_amount;type:decimal(10, 2);comment:预估金额（以最高金额计算，后期据情况退回部分金额）【系统计算】;not null;" json:"estimated_amount"`                                 // 预估金额（以最高金额计算，后期据情况退回部分金额）【系统计算】
	CouponId              int64     `gorm:"column:coupon_id;type:bigint;comment:优惠券ID;default:NULL;" json:"coupon_id"`                                                                             // 优惠券ID
	CouponAmount          float64   `gorm:"column:coupon_amount;type:decimal(10, 2);comment:优惠金额;default:0.00;" json:"coupon_amount"`                                                              // 优惠金额
	ActualAmount          float64   `gorm:"column:actual_amount;type:decimal(10, 2);comment:实际金额（不包括高速费）;default:NULL;" json:"actual_amount"`                                                      // 实际金额（不包括高速费）
	TollFee               float64   `gorm:"column:toll_fee;type:decimal(10, 2);comment:高速费（司机输入、可协商、可直接转司机）;default:0.00;" json:"toll_fee"`                                                        // 高速费（司机输入、可协商、可直接转司机）
	PaymentMethod         string    `gorm:"column:payment_method;type:varchar(10);comment:支付方式：1支付宝 2微信 3银联 4余额;default:NULL;" json:"payment_method"`                                              // 支付方式：1支付宝 2微信 3银联 4余额
	PaymentStatus         string    `gorm:"column:payment_status;type:varchar(10);comment:支付状态：1-待支付 2已支付;default:1;" json:"payment_status"`                                                       // 支付状态：1-待支付 2已支付
	PaymentTime           time.Time `gorm:"column:payment_time;type:datetime;comment:支付时间;default:NULL;" json:"payment_time"`                                                                      // 支付时间
	OrderStatus           string    `gorm:"column:order_status;type:varchar(10);comment:订单状态：1-待接单 2-已接单 3-待司机到达 4-司机已到达 5-进行中 6-已结束 7-已完成 8-已取消;not null;default:1;" json:"order_status"`         // 订单状态：1-待接单 2-已接单 3-待司机到达 4-司机已到达 5-进行中 6-已结束 7-已完成 8-已取消
	CancelTime            time.Time `gorm:"column:cancel_time;type:datetime;comment:取消时间;default:NULL;" json:"cancel_time"`                                                                        // 取消时间
	CancelReason          string    `gorm:"column:cancel_reason;type:varchar(255);comment:取消原因;default:NULL;" json:"cancel_reason"`                                                                // 取消原因
	CreatedAt             time.Time `gorm:"column:created_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"updated_at"`
	DeletedAt             time.Time `gorm:"column:deleted_at;type:datetime(3);default:NULL;" json:"deleted_at"`
}

func (so *ShunOrder) TableName() string {
	return "shun_order"
}
func (so *ShunOrder) CreateOrder(db *gorm.DB) error {
	return db.Create(&so).Error
}
func (so *ShunOrder) Editor(db *gorm.DB) error {
	return db.Updates(&so).Error
}
func (so *ShunOrder) GetOrderByCode(db *gorm.DB, code string) error {
	return db.Where("order_code = ?", code).First(&so).Error
}
func (so *ShunOrder) GetOrderById(db *gorm.DB, id int64) error {
	return db.Where("id = ?", id).First(&so).Error
}

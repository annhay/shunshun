package model

import "time"

type ShunCoupon struct { //优惠券信息表
	Id            uint64    `gorm:"column:id;type:bigint UNSIGNED;primaryKey;not null;" json:"id"`
	SnCode        string    `gorm:"column:sn_code;type:varchar(20);comment:优惠券编号;not null;" json:"sn_code"`                                    // 优惠券编号
	Name          string    `gorm:"column:name;type:varchar(20);comment:优惠券名称;not null;" json:"name"`                                          // 优惠券名称
	CouponType    string    `gorm:"column:coupon_type;type:varchar(10);comment:优惠券类型：1-满减券 2-折扣券 3-无门槛券;not null;" json:"coupon_type"`         // 优惠券类型：1-满减券 2-折扣券 3-无门槛券
	IsLimitPut    string    `gorm:"column:is_limit_put;type:varchar(10);comment:限制发放数量:1-是 2-否;not null;default:2;" json:"is_limit_put"`       // 限制发放数量:1-是 2-否
	PutCount      int32     `gorm:"column:put_count;type:int;comment:发放数量;default:NULL;" json:"put_count"`                                     // 发放数量
	UseCouponType string    `gorm:"column:use_coupon_type;type:varchar(10);comment:用券类型:1固定日期 2领券当日起 3领券次日起;not null;" json:"use_coupon_type"` // 用券类型:1固定日期 2领券当日起 3领券次日起
	StartTime     time.Time `gorm:"column:start_time;type:datetime;comment:领券开始时间;default:NULL;" json:"start_time"`                            // 领券开始时间
	EndTime       time.Time `gorm:"column:end_time;type:datetime;comment:领券结束时间;default:NULL;" json:"end_time"`                                // 领券结束时间
	DayNum        int32     `gorm:"column:day_num;type:int;comment:可用期限;default:NULL;" json:"day_num"`                                         // 可用期限
	Money         float64   `gorm:"column:money;type:decimal(10, 2);comment:优惠额度;default:NULL;" json:"money"`                                  // 优惠额度
	SpendMoney    float64   `gorm:"column:spend_money;type:decimal(10, 2);comment:满多少;default:NULL;" json:"spend_money"`                       // 满多少
	Discount      int32     `gorm:"column:discount;type:int;comment:折扣(需要换算);default:NULL;" json:"discount"`                                   // 折扣(需要换算)
	TotalMoney    float64   `gorm:"column:total_money;type:decimal(10, 2);comment:最高优惠金额(只有折扣券);default:NULL;" json:"total_money"`             // 最高优惠金额(只有折扣券)
	GetType       string    `gorm:"column:get_type;type:varchar(10);comment:领取类型：1-不限制领取次数 2-限制领取 3-每天限制领取;not null;" json:"get_type"`         // 领取类型：1-不限制领取次数 2-限制领取 3-每天限制领取
	LimitCount    int32     `gorm:"column:limit_count;type:int;comment:限制次数;default:NULL;" json:"limit_count"`                                 // 限制次数
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"updated_at"`
	DeletedAt     time.Time `gorm:"column:deleted_at;type:datetime(3);default:NULL;" json:"deleted_at"`
}

func (sc *ShunCoupon) TableName() string {
	return "shun_coupon"
}

type ShunUseCoupon struct { //优惠使用情况表
	Id             uint64    `gorm:"column:id;type:bigint UNSIGNED;primaryKey;not null;" json:"id"`
	UserId         int64     `gorm:"column:user_id;type:bigint;comment:用户ID;default:NULL;" json:"user_id"`                         // 用户ID
	CouponId       int32     `gorm:"column:coupon_id;type:int;comment:优惠券ID;default:NULL;" json:"coupon_id"`                       // 优惠券ID
	UseStatus      int8      `gorm:"column:use_status;type:tinyint;comment:使用状态:1未使用 2已使用 3已过期 4已作废;default:1;" json:"use_status"` // 使用状态:1未使用 2已使用 3已过期 4已作废
	UseTime        time.Time `gorm:"column:use_time;type:datetime;comment:使用时间;default:NULL;" json:"use_time"`                     // 使用时间
	ExpirationTime time.Time `gorm:"column:expiration_time;type:datetime;comment:过期时间;default:NULL;" json:"expiration_time"`       // 过期时间
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"updated_at"`
	DeletedAt      time.Time `gorm:"column:deleted_at;type:datetime(3);default:NULL;" json:"deleted_at"`
}

func (su *ShunUseCoupon) TableName() string {
	return "shun_use_coupon"
}

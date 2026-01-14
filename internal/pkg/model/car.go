package model

import (
	"time"

	"gorm.io/gorm"
)

type ShunCar struct { //车辆信息
	Id                  uint64    `gorm:"column:id;type:bigint UNSIGNED;primaryKey;not null;" json:"id"`
	DriverId            uint64    `gorm:"column:driver_id;type:bigint UNSIGNED;comment:司机ID;not null;" json:"driver_id"`                        // 司机ID
	VehicleNo           string    `gorm:"column:vehicle_no;type:varchar(20);comment:车牌号;not null;" json:"vehicle_no"`                           // 车牌号
	VehicleType         string    `gorm:"column:vehicle_type;type:varchar(10);comment:车辆类型;not null;" json:"vehicle_type"`                      // 车辆类型
	VehicleBrand        string    `gorm:"column:vehicle_brand;type:varchar(32);comment:车辆品牌;not null;" json:"vehicle_brand"`                    // 车辆品牌
	VehicleModel        string    `gorm:"column:vehicle_model;type:varchar(64);comment:车辆车型;not null;" json:"vehicle_model"`                    // 车辆车型
	VehicleColor        string    `gorm:"column:vehicle_color;type:varchar(20);comment:车辆颜色;not null;" json:"vehicle_color"`                    // 车辆颜色
	Vin                 string    `gorm:"column:vin;type:varchar(32);comment:车辆识别码;not null;" json:"vin"`                                       // 车辆识别码
	EngineNo            string    `gorm:"column:engine_no;type:varchar(32);comment:发动机号;default:NULL;" json:"engine_no"`                        // 发动机号
	RegisterDate        time.Time `gorm:"column:register_date;type:datetime;comment:注册日期;not null;" json:"register_date"`                       // 注册日期
	LicenseNo           string    `gorm:"column:license_no;type:varchar(64);comment:行驶证编号;not null;" json:"license_no"`                         // 行驶证编号
	LicenseExpireDate   time.Time `gorm:"column:license_expire_date;type:datetime;comment:行驶证过期日期;not null;" json:"license_expire_date"`        // 行驶证过期日期
	InsuranceExpireDate time.Time `gorm:"column:Insurance_expire_date;type:datetime;comment:保险到期日期;not null;" json:"Insurance_expire_date"`     // 保险到期日期
	Status              string    `gorm:"column:status;type:varchar(10);comment:车辆状态：0待审核 1正常 2审核驳回 3停用 4注销;not null;default:0;" json:"status"` // 车辆状态：0待审核 1正常 2审核驳回 3停用 4注销
	RejectReason        string    `gorm:"column:reject_reason;type:varchar(255);comment:驳回原因;default:NULL;" json:"reject_reason"`               // 驳回原因
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"updated_at"`
	DeletedAt           time.Time `gorm:"column:deleted_at;type:datetime(3);default:NULL;" json:"deleted_at"`
}

func (sc *ShunCar) TableName() string {
	return "shun_car"
}
func (sc *ShunCar) CreateCar(db *gorm.DB) error {
	return db.Create(&sc).Error
}
func (sc *ShunCar) GetCarsByDriverId(db *gorm.DB, driverId uint64) ([]ShunCar, error) {
	var cars []ShunCar
	err := db.Where("driver_id = ?", driverId).Find(&cars).Error
	return cars, err
}
func (sc *ShunCar) UpdateCar(db *gorm.DB) error {
	return db.Save(&sc).Error
}

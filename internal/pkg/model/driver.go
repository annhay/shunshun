package model

import (
	"time"

	"gorm.io/gorm"
)

// ShunDriver 司机信息模型
// 存储司机的基本信息、身份信息、驾驶证信息等

type ShunDriver struct { //司机信息表
	Id                      uint64    `gorm:"column:id;type:bigint UNSIGNED;primaryKey;not null;" json:"id"`
	UserId                  uint64    `gorm:"column:user_id;type:bigint UNSIGNED;comment:用户ID;not null;" json:"user_id"`                                       // 用户 ID
	DriverNo                string    `gorm:"column:driver_no;type:varchar(32);comment:司机编号;not null;" json:"driver_no"`                                       // 司机编号
	RealName                string    `gorm:"column:real_name;type:varchar(64);comment:真实姓名;not null;" json:"real_name"`                                       // 真实姓名
	IdCardNo                string    `gorm:"column:id_card_no;type:varchar(50);comment:身份证号;not null;" json:"id_card_no"`                                     // 身份证号
	IdCardFrontUrl          string    `gorm:"column:id_card_front_url;type:varchar(255);comment:身份证正面URL;not null;" json:"id_card_front_url"`                  // 身份证正面 URL
	IdCardBackUrl           string    `gorm:"column:id_card_back_url;type:varchar(255);comment:身份证反面URL;not null;" json:"id_card_back_url"`                    // 身份证反面 URL
	IdCardExpireTime        time.Time `gorm:"column:id_card_expire_time;type:datetime;comment:身份证有效期;not null;" json:"id_card_expire_time"`                    // 身份证有效期
	DriverLicenseNo         string    `gorm:"column:driver_license_no;type:varchar(32);comment:驾驶证号;not null;" json:"driver_license_no"`                       // 驾驶证号
	DriverLicenseUrl        string    `gorm:"column:driver_license_url;type:varchar(255);comment:驾驶证URL;not null;" json:"driver_license_url"`                  // 驾驶证 URL
	DriverLicenseGetTime    time.Time `gorm:"column:driver_license_get_time;type:datetime;comment:驾驶证领证时间;not null;" json:"driver_license_get_time"`           // 驾驶证领证时间
	DriverLicenseExpireTime time.Time `gorm:"column:driver_license_expire_time;type:datetime;comment:驾驶证有效期;not null;" json:"driver_license_expire_time"`      // 驾驶证有效期
	DrivingAge              uint8     `gorm:"column:driving_age;type:tinyint UNSIGNED;comment:驾龄(年);default:NULL;" json:"driving_age"`                         // 驾龄(年)
	HealthCertUrl           string    `gorm:"column:health_cert_url;type:varchar(255);comment:体检报告URL;default:NULL;" json:"health_cert_url"`                   // 体检报告 URL
	ResidencePermitUrl      string    `gorm:"column:residence_permit_url;type:varchar(255);comment:居住证URL;default:NULL;" json:"residence_permit_url"`          // 居住证 URL
	CityCode                string    `gorm:"column:city_code;type:varchar(6);comment:城市编码;default:NULL;" json:"city_code"`                                    // 城市编码
	AuditStatus             string    `gorm:"column:audit_status;type:varchar(10);comment:审核状态  1-待审核 2-通过 3-驳回 4-禁用;not null;default:1;" json:"audit_status"` // 审核状态  1-待审核 2-通过 3-驳回 4-禁用
	CreatedAt               time.Time `gorm:"column:created_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"created_at"`
	UpdatedAt               time.Time `gorm:"column:updated_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"updated_at"`
	DeletedAt               time.Time `gorm:"column:deleted_at;type:datetime(3);default:NULL;" json:"deleted_at"`
}

// TableName 指定表名
//
// 返回值:
//   - string: 表名
func (sd *ShunDriver) TableName() string {
	return "shun_driver"
}

// CreateDriver 创建司机信息
//
// 参数:
//   - db *gorm.DB: 数据库连接
//
// 返回值:
//   - error: 错误信息
func (sd *ShunDriver) CreateDriver(db *gorm.DB) error {
	return db.Create(&sd).Error
}

// GetDriverByUserId 根据用户ID获取司机信息
//
// 参数:
//   - db *gorm.DB: 数据库连接
//   - userId int64: 用户ID
//
// 返回值:
//   - error: 错误信息
func (sd *ShunDriver) GetDriverByUserId(db *gorm.DB, userId int64) error {
	return db.Where("user_id = ?", userId).First(&sd).Error
}

// Editor 更新司机信息
//
// 参数:
//   - db *gorm.DB: 数据库连接
//
// 返回值:
//   - error: 错误信息
func (sd *ShunDriver) Editor(db *gorm.DB) error {
	return db.Updates(&sd).Error
}

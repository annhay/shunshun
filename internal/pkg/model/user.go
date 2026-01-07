package model

import (
	"time"

	"gorm.io/gorm"
)

type ShunUser struct { //用户
	Id             uint64    `gorm:"column:id;type:bigint UNSIGNED;primaryKey;not null;" json:"id"`
	Phone          string    `gorm:"column:phone;type:char(11);comment:手机号;not null;" json:"phone"`                                                // 手机号
	Password       string    `gorm:"column:password;type:varchar(50);comment:密码;not null;" json:"password"`                                        // 密码
	Cover          string    `gorm:"column:cover;type:varchar(255);comment:头像;default:NULL;" json:"cover"`                                         // 头像
	Nickname       string    `gorm:"column:nickname;type:varchar(50);comment:昵称;default:NULL;" json:"nickname"`                                    // 昵称
	Sex            string    `gorm:"column:sex;type:varchar(10);comment:性别:1男 2女;default:NULL;" json:"sex"`                                        // 性别:1男 2女
	RealName       string    `gorm:"column:real_name;type:varchar(50);comment:真实姓名;default:NULL;" json:"real_name"`                                // 真实姓名
	IdCard         string    `gorm:"column:id_card;type:char(18);comment:身份证号;default:NULL;" json:"id_card"`                                       // 身份证号
	BirthdayTime   time.Time `gorm:"column:birthday_time;type:datetime;comment:出生日期;default:NULL;" json:"birthday_time"`                           // 出生日期
	SchoolName     string    `gorm:"column:school_name;type:varchar(100);comment:学校名称;default:NULL;" json:"school_name"`                           // 学校名称
	StudentId      string    `gorm:"column:student_id;type:varchar(50);comment:学号;default:NULL;" json:"student_id"`                                // 学号
	EnrollmentYear time.Time `gorm:"column:enrollment_year;type:datetime;comment:入学年份;default:NULL;" json:"enrollment_year"`                       // 入学年份
	StudentIdPhoto string    `gorm:"column:student_id_photo;type:varchar(255);comment:学生证照片;default:NULL;" json:"student_id_photo"`                // 学生证照片
	AuthStatus     string    `gorm:"column:auth_status;type:varchar(10);comment:认证状态:0未认证 1审核中 2已认证 3认证失败;not null;default:0;" json:"auth_status"` // 认证状态:0未认证 1审核中 2已认证 3认证失败
	ExpirationTime time.Time `gorm:"column:expiration_time;type:datetime;comment:学生优惠过期时间=学生证过期时间+1年;default:NULL;" json:"expiration_time"`        // 学生优惠过期时间=学生证过期时间+1年
	LastLoginTime  time.Time `gorm:"column:last_login_time;type:datetime;comment:最后登录时间;default:NULL;" json:"last_login_time"`                     // 最后登录时间
	Status         string    `gorm:"column:status;type:varchar(10);comment:状态:1正常 2封禁 3注销;not null;default:1;" json:"status"`                      // 状态:1正常 2封禁 3注销
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"updated_at"`
	DeletedAt      time.Time `gorm:"column:deleted_at;type:datetime(3);default:NULL;" json:"deleted_at"`
}

func (su *ShunUser) TableName() string {
	return "shun_user"
}
func (su *ShunUser) GetUserByPhone(db *gorm.DB, phone string) error {
	return db.Where("phone = ?", phone).First(&su).Error
}
func (su *ShunUser) CreateUser(db *gorm.DB) error {
	return db.Create(&su).Error
}
func (su *ShunUser) Editor(db *gorm.DB) error {
	return db.Updates(&su).Error
}
func (su *ShunUser) GetUserById(db *gorm.DB, id int) error {
	return db.Where("id = ?", id).First(&su).Error
}

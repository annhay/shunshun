package model

type ShunUser struct { //用户表
	Id       uint64 `gorm:"column:id;type:bigint UNSIGNED;primaryKey;not null;" json:"id"`
	Phone    string `gorm:"column:phone;type:char(11);comment:手机号;not null;" json:"phone"`                       // 手机号
	Password string `gorm:"column:password;type:varchar(50);comment:密码;not null;" json:"password"`                // 密码
	Cover    string `gorm:"column:cover;type:varchar(255);comment:头像;default:NULL;" json:"cover"`                 // 头像
	Nickname string `gorm:"column:nickname;type:varchar(50);comment:昵称;default:NULL;" json:"nickname"`            // 昵称
	Sex      string `gorm:"column:sex;type:varchar(10);comment:性别:1男 2女;default:NULL;" json:"sex"`              // 性别:1男 2女
	RealName string `gorm:"column:real_name;type:varchar(50);comment:真实姓名;default:NULL;" json:"real_name"`      // 真实姓名
	IdCard   string `gorm:"column:id_card;type:char(18);comment:身份证号;default:NULL;" json:"id_card"`             // 身份证号
	Status   string `gorm:"column:status;type:varchar(10);comment:状态:1正常 2封禁 3注销;default:1;" json:"status"` // 状态:1正常 2封禁 3注销
}

func (su *ShunUser) TableName() string {
	return "shun_user"
}

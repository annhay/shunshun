package request

// SendTextMessage 绑定 JSON
type SendTextMessage struct {
	Phone string `form:"phone" json:"phone" xml:"phone"  binding:"required"`
}

// Register 绑定 JSON
type Register struct {
	Phone            string `form:"phone" json:"phone" xml:"phone"  binding:"required"`
	VerificationCode string `form:"verification_code" json:"verification_code" xml:"verification_code"  binding:"required"`
	Password         string `form:"password" json:"password" xml:"password"  binding:"required"`
}

// Login 绑定 JSON
type Login struct {
	Phone            string `form:"phone" json:"phone" xml:"phone"  binding:"required"`
	Password         string `form:"password" json:"password" xml:"password"  binding:"required"`
	VerificationCode string `form:"verification_code" json:"verification_code" xml:"verification_code"  binding:"required"`
}

// ForgotPassword 绑定 JSON
type ForgotPassword struct {
	Phone            string `form:"phone" json:"phone" xml:"phone"  binding:"required"`
	VerificationCode string `form:"verification_code" json:"verification_code" xml:"verification_code"  binding:"required"`
	Password         string `form:"password" json:"password" xml:"password"  binding:"required"`
	ConfirmPassword  string `form:"confirm_password" json:"confirm_password" xml:"confirm_password"  binding:"required"`
}

// CompleteInformation 绑定 JSON
type CompleteInformation struct {
	Cover        string `form:"cover" json:"cover" xml:"cover"`
	Nickname     string `form:"nickname" json:"nickname" xml:"nickname"`
	Sex          string `form:"sex" json:"sex" xml:"sex"`
	RealName     string `form:"real_name" json:"real_name" xml:"real_name"`
	IdCard       string `form:"id_card" json:"id_card" xml:"id_card"`
	BirthdayTime string `form:"birthday_time" json:"birthday_time" xml:"birthday_time"`
}

// StudentVerification 绑定 JSON
type StudentVerification struct {
	SchoolName     string `form:"school_name" json:"school_name" xml:"school_name"  binding:"required"`
	StudentId      string `form:"student_id" json:"student_id" xml:"student_id"  binding:"required"`
	EnrollmentYear string `form:"enrollment_year" json:"enrollment_year" xml:"enrollment_year"  binding:"required"`
	StudentIdPhoto string `form:"student_id_photo" json:"student_id_photo" xml:"student_id_photo"  binding:"required"`
}

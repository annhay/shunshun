package request

// SendTextMessage 绑定 JSON
type SendTextMessage struct {
	Phone string `form:"phone" json:"phone" xml:"phone"  binding:"required"`
}

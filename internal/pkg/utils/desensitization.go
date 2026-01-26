package utils

import "strings"

// 敏感信息脱敏工具
// 用于对用户敏感信息进行脱敏处理，保护用户隐私

// PhoneDesensitization 手机号脱敏
// 
// 参数:
//   - phone: 手机号字符串
// 
// 返回值:
//   - string: 脱敏后的手机号，显示前3位和后4位，中间的4位用星号代替
// 
// 使用场景:
//   - 在用户个人中心、订单详情等页面显示用户手机号时使用
//   - 确保手机号部分可见，同时保护用户隐私
func PhoneDesensitization(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}

// IDCardDesensitization 身份证号脱敏
// 
// 参数:
//   - idCard: 身份证号字符串
// 
// 返回值:
//   - string: 脱敏后的身份证号，显示前6位和后4位，中间的用星号代替
// 
// 使用场景:
//   - 在需要显示身份证号的场景中使用
//   - 保护用户身份证号隐私，同时保留部分可识别信息
func IDCardDesensitization(idCard string) string {
	length := len(idCard)
	if length != 15 && length != 18 {
		return idCard
	}
	return idCard[:6] + strings.Repeat("*", length-10) + idCard[length-4:]
}

// StudentIDDesensitization 学生证号脱敏
// 
// 参数:
//   - studentID: 学生证号字符串
// 
// 返回值:
//   - string: 脱敏后的学生证号，显示前2位和后3位，中间的用星号代替
// 
// 使用场景:
//   - 在学生认证、优惠信息等场景中显示学生证号时使用
//   - 保护学生隐私，同时保留部分可识别信息
func StudentIDDesensitization(studentID string) string {
	length := len(studentID)
	if length < 6 {
		return studentID
	}

	return studentID[:2] + strings.Repeat("*", length-5) + studentID[length-3:]
}

// NameDesensitization 姓名脱敏
// 
// 参数:
//   - name: 姓名字符串
// 
// 返回值:
//   - string: 脱敏后的姓名，显示姓，名用星号代替
// 
// 使用场景:
//   - 在需要显示用户姓名的场景中使用
//   - 保护用户姓名隐私，同时保留姓氏信息
func NameDesensitization(name string) string {
	length := len(name)
	if length <= 1 {
		return name
	}
	return string(name[0]) + strings.Repeat("*", length-1)
}

package utils

import "strings"

// SensitiveInfoDesensitization 敏感信息脱敏工具

// PhoneDesensitization 手机号脱敏
// 显示前3位和后4位，中间的4位用星号代替
func PhoneDesensitization(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}

// IDCardDesensitization 身份证号脱敏
// 显示前6位和后4位，中间的8位用星号代替
func IDCardDesensitization(idCard string) string {
	length := len(idCard)
	if length != 15 && length != 18 {
		return idCard
	}
	return idCard[:6] + strings.Repeat("*", length-10) + idCard[length-4:]
}

// StudentIDDesensitization 学生证号脱敏
// 显示前2位和后3位，中间的用星号代替
func StudentIDDesensitization(studentID string) string {
	length := len(studentID)
	if length < 6 {
		return studentID
	}

	return studentID[:2] + strings.Repeat("*", length-5) + studentID[length-3:]
}

// NameDesensitization 姓名脱敏
// 显示姓，名用星号代替
func NameDesensitization(name string) string {
	length := len(name)
	if length <= 1 {
		return name
	}
	return string(name[0]) + strings.Repeat("*", length-1)
}

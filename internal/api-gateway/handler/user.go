package handler

import (
	"net/http"
	"shunshun/internal/api-gateway/consts"
	"shunshun/internal/api-gateway/middleware"
	"shunshun/internal/api-gateway/request"
	"shunshun/internal/pkg/global"
	"shunshun/internal/proto"

	"github.com/gin-gonic/gin"
)

// SendTextMessage 短信发送
func SendTextMessage(c *gin.Context) {
	var form request.SendTextMessage
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := global.UserClient.SendTextMessage(c, &proto.SendTextMessageReq{Phone: form.Phone})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "短信发送成功"})
}

// Register 注册
func Register(c *gin.Context) {
	var form request.Register
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	register, err := global.UserClient.Register(c, &proto.RegisterReq{
		Phone:            form.Phone,
		VerificationCode: form.VerificationCode,
		Password:         form.Password,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	token := middleware.CreateToken(uint64(register.Id), consts.JwtKey)
	c.JSON(http.StatusOK, gin.H{"message": "注册成功", "token": token})
}

// Login 登录
func Login(c *gin.Context) {
	var form request.Login
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	login, err := global.UserClient.Login(c, &proto.LoginReq{
		Phone:            form.Phone,
		Password:         form.Password,
		VerificationCode: form.VerificationCode,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	token := middleware.CreateToken(uint64(login.Id), consts.JwtKey)
	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "token": token})
}

// ForgotPassword 忘记密码
func ForgotPassword(c *gin.Context) {
	var form request.ForgotPassword
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := global.UserClient.ForgotPassword(c, &proto.ForgotPasswordReq{
		Phone:            form.Phone,
		VerificationCode: form.VerificationCode,
		Password:         form.Password,
		ConfirmPassword:  form.ConfirmPassword,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "修改成功"})
}

// CompleteInformation 完善信息
func CompleteInformation(c *gin.Context) {
	var form request.CompleteInformation
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := global.UserClient.CompleteInformation(c, &proto.CompleteInformationReq{
		UserId:       int64(c.GetUint("userId")),
		Cover:        form.Cover,
		Nickname:     form.Nickname,
		Sex:          form.Sex,
		RealName:     form.RealName,
		IdCard:       form.IdCard,
		BirthdayTime: form.BirthdayTime,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "完善信息成功"})
}

// StudentVerification 学生认证
func StudentVerification(c *gin.Context) {
	var form request.StudentVerification
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := global.UserClient.StudentVerification(c, &proto.StudentVerificationReq{
		UserId:         int64(c.GetUint("userId")),
		SchoolName:     form.SchoolName,
		StudentId:      form.StudentId,
		EnrollmentYear: form.EnrollmentYear,
		StudentIdPhoto: form.StudentIdPhoto,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "认证成功"})
}

// PersonalCenter 个人中心
func PersonalCenter(c *gin.Context) {
	center, err := global.UserClient.PersonalCenter(c, &proto.PersonalCenterReq{
		UserId: int64(c.GetUint("userId")),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": center})
}

// Logout 账号注销
func Logout(c *gin.Context) {
	_, err := global.UserClient.Logout(c, &proto.LogoutReq{
		UserId: int64(c.GetUint("userId")),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "账号已注销"})
}

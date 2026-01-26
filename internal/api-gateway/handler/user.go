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

// SendTextMessage 发送短信验证码
// 
// 参数:
//   - c *gin.Context: Gin上下文
// 
// 处理逻辑:
//   1. 绑定请求参数
//   2. 调用用户服务发送短信
//   3. 返回发送结果
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

// Register 用户注册
// 
// 参数:
//   - c *gin.Context: Gin上下文
// 
// 处理逻辑:
//   1. 绑定请求参数
//   2. 调用用户服务进行注册
//   3. 生成JWT令牌
//   4. 返回注册结果和令牌
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

// Login 用户登录
// 
// 参数:
//   - c *gin.Context: Gin上下文
// 
// 处理逻辑:
//   1. 绑定请求参数
//   2. 调用用户服务进行登录
//   3. 生成JWT令牌
//   4. 返回登录结果和令牌
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
// 
// 参数:
//   - c *gin.Context: Gin上下文
// 
// 处理逻辑:
//   1. 绑定请求参数
//   2. 调用用户服务进行密码重置
//   3. 返回重置结果
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

// CompleteInformation 完善用户信息
// 
// 参数:
//   - c *gin.Context: Gin上下文
// 
// 处理逻辑:
//   1. 绑定请求参数
//   2. 从上下文获取用户ID
//   3. 调用用户服务完善信息
//   4. 返回完善结果
func CompleteInformation(c *gin.Context) {
	var form request.CompleteInformation
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := global.UserClient.CompleteInformation(c, &proto.CompleteInformationReq{
		UserId:       int64(c.GetUint64("userId")),
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
// 
// 参数:
//   - c *gin.Context: Gin上下文
// 
// 处理逻辑:
//   1. 绑定请求参数
//   2. 从上下文获取用户ID
//   3. 调用用户服务进行学生认证
//   4. 返回认证结果
func StudentVerification(c *gin.Context) {
	var form request.StudentVerification
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := global.UserClient.StudentVerification(c, &proto.StudentVerificationReq{
		UserId:         int64(c.GetUint64("userId")),
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
// 
// 参数:
//   - c *gin.Context: Gin上下文
// 
// 处理逻辑:
//   1. 从上下文获取用户ID
//   2. 调用用户服务获取个人中心信息
//   3. 返回个人中心信息
func PersonalCenter(c *gin.Context) {
	center, err := global.UserClient.PersonalCenter(c, &proto.PersonalCenterReq{
		UserId: int64(c.GetUint64("userId")),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": center})
}

// Logout 账号注销
// 
// 参数:
//   - c *gin.Context: Gin上下文
// 
// 处理逻辑:
//   1. 从上下文获取用户ID
//   2. 调用用户服务进行账号注销
//   3. 返回注销结果
func Logout(c *gin.Context) {
	_, err := global.UserClient.Logout(c, &proto.LogoutReq{
		UserId: int64(c.GetUint64("userId")),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "账号已注销"})
}

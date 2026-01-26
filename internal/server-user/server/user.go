package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"shunshun/internal/pkg/global"
	"shunshun/internal/pkg/model"
	"shunshun/internal/pkg/utils"
	"shunshun/internal/proto"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Server 用户服务结构体
// 实现了proto.UnimplementedUserServer接口

type Server struct {
	proto.UnimplementedUserServer
}

// SendTextMessage 短信发送
//
// 参数:
//   - _ context.Context: 上下文（未使用）
//   - in *proto.SendTextMessageReq: 短信发送请求，包含手机号
//
// 返回值:
//   - *proto.SendTextMessageResp: 短信发送响应
//   - error: 错误信息
func (s *Server) SendTextMessage(_ context.Context, in *proto.SendTextMessageReq) (*proto.SendTextMessageResp, error) {
	// 1分钟发送次数Key
	minuteSendCountKey := fmt.Sprintf("%s:minute_send_text-message_count", in.Phone)
	// 获取发送次数
	count, _ := global.Rdb.Get(context.Background(), minuteSendCountKey).Int()
	// >=1 发送拦截
	if count >= 1 {
		return nil, errors.New("一分钟内只能发送1次短信")
	}

	// 1小时发送次数key
	hourSendCountKey := fmt.Sprintf("%s:hour_send_text-message_count", in.Phone)
	// 获取发送次数
	counts, _ := global.Rdb.Get(context.Background(), hourSendCountKey).Int()
	// >=5发送拦截
	if counts >= 5 {
		return nil, errors.New("一小时内只能发送5次短信")
	}

	// 随机数生成验证码
	code := rand.Intn(9000) + 1000
	// 互亿第三方短信平台验证
	utils.HuYi(in.Phone, strconv.Itoa(code))

	// 设置发送短信验证码的键
	sendKey := fmt.Sprintf("%s:send_text-message", in.Phone)
	// 存入redis,5分钟过期时间
	global.Rdb.Set(context.Background(), sendKey, code, time.Minute*5)

	// 记录一分钟内的发送次数
	incr := global.Rdb.Incr(context.Background(), minuteSendCountKey)
	if incr.Val() == 1 {
		// 设置1分钟过期时间
		global.Rdb.Expire(context.Background(), minuteSendCountKey, time.Minute*1)
	}

	// 记录1小时内的发送次数
	hourIncr := global.Rdb.Incr(context.Background(), hourSendCountKey)
	if hourIncr.Val() == 1 {
		// 设置1小时的过期时间
		global.Rdb.Expire(context.Background(), hourSendCountKey, time.Hour*1)
	}

	return &proto.SendTextMessageResp{}, nil
}

// Register 注册(自动登录)
//
// 参数:
//   - _ context.Context: 上下文（未使用）
//   - in *proto.RegisterReq: 注册请求，包含手机号、密码、验证码
//
// 返回值:
//   - *proto.RegisterResp: 注册响应，包含用户ID
//   - error: 错误信息
func (s *Server) Register(_ context.Context, in *proto.RegisterReq) (*proto.RegisterResp, error) {
	// 账号查询，防止重复账号注册
	var user model.ShunUser
	if err := user.GetUserByPhone(global.DB, utils.EnPwdCode([]byte(in.Phone))); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账号已存在")
		}
	}

	// 从缓存中获取短信验证码
	sendKey := fmt.Sprintf("%s:send_text-message", in.Phone)
	result, err := global.Rdb.Get(context.Background(), sendKey).Result()
	if err != nil {
		return nil, err
	}

	// 验证码判断
	if in.VerificationCode != result {
		return nil, errors.New("短信验证码错误,请重新输入")
	}
	// 设置用户信息
	user.Phone = utils.EnPwdCode([]byte(in.Phone)) //储存加密后的数据
	user.Password = utils.Md5(in.Password)
	// 设置默认头像
	user.Cover = "https://gd-hbimg.huaban.com/9980f53d02800ebad0472f0fe1a6eed9a1ba699ec27f-mvUjPb_fw658"
	// 创建用户
	if err := user.CreateUser(global.DB); err != nil {
		return nil, err
	}
	// 设置昵称，格式为：用户+注册时间+用户ID
	user.Nickname = "用户" + time.Now().Format("20060102150405") + strconv.FormatUint(user.Id, 10)
	// 注册成功后自动登录并更新最后登录时间
	user.LastLoginTime = time.Now()

	// 更新用户信息
	if err := user.Editor(global.DB); err != nil {
		return nil, err
	}

	// 最后删除缓存中的验证码，防止重复使用
	global.Rdb.Del(context.Background(), sendKey)

	return &proto.RegisterResp{Id: int64(user.Id)}, nil
}

// Login 登录
//
// 参数:
//   - _ context.Context: 上下文（未使用）
//   - in *proto.LoginReq: 登录请求，包含手机号、密码/验证码
//
// 返回值:
//   - *proto.LoginResp: 登录响应，包含用户ID
//   - error: 错误信息
func (s *Server) Login(_ context.Context, in *proto.LoginReq) (*proto.LoginResp, error) {
	// 账号查询判断
	var user model.ShunUser
	if err := user.GetUserByPhone(global.DB, utils.EnPwdCode([]byte(in.Phone))); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账号不存在")
		}
	}
	//判断账号状态
	if user.Status == "2" { //封禁
		return nil, errors.New("账号已被封禁")
	}
	if user.Status == "3" { //注销
		return nil, errors.New("账号已被注销")
	}
	// 判断登录方式
	if in.Password == "" && in.VerificationCode != "" { // 短信验证码登录
		// 从缓存中获取短信验证码
		sendKey := fmt.Sprintf("%s:send_text-message", in.Phone)
		result, err := global.Rdb.Get(context.Background(), sendKey).Result()
		if err != nil {
			return nil, err
		}

		// 验证码判断
		if in.VerificationCode != result {
			return nil, errors.New("短信验证码错误,请重新输入")
		}

		// 最后删除缓存中的验证码，防止重复使用
		global.Rdb.Del(context.Background(), sendKey)
	} else { // 密码登录
		// 设置错误次数键
		banKey := fmt.Sprintf("%s:login_ban", in.Phone)
		// 获取错误次数
		count, err := global.Rdb.Get(context.Background(), banKey).Int()
		if err != nil {
			return nil, err
		}

		// 错误次数>=3登录拦截
		if count >= 3 {
			return nil, errors.New("登录次数过多,请2小时后在登录")
		}

		// 密码错误
		if user.Password != utils.Md5(in.Password) {
			// 记录错误次数
			incr := global.Rdb.Incr(context.Background(), banKey)
			if incr.Val() == 1 {
				// 设置过期时间2小时
				global.Rdb.Expire(context.Background(), banKey, time.Hour*2)
			}
			return nil, errors.New("密码错误,请重新输入")
		}

		// 登录成功,删除错误次数键
		global.Rdb.Del(context.Background(), banKey)
	}

	// 更新最后登录时间
	user.LastLoginTime = time.Now()
	if err := user.Editor(global.DB); err != nil {
		return nil, err
	}

	return &proto.LoginResp{Id: int64(user.Id)}, nil
}

// ForgotPassword 忘记密码
//
// 参数:
//   - _ context.Context: 上下文（未使用）
//   - in *proto.ForgotPasswordReq: 忘记密码请求，包含手机号、新密码、确认密码、验证码
//
// 返回值:
//   - *proto.ForgotPasswordResp: 忘记密码响应
//   - error: 错误信息
func (s *Server) ForgotPassword(_ context.Context, in *proto.ForgotPasswordReq) (*proto.ForgotPasswordResp, error) {
	// 账号查询判断
	var user model.ShunUser
	if err := user.GetUserByPhone(global.DB, utils.EnPwdCode([]byte(in.Phone))); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账号不存在")
		}
	}

	// 从缓存中获取短信验证码
	sendKey := fmt.Sprintf("%s:send_text-message", in.Phone)
	result, err := global.Rdb.Get(context.Background(), sendKey).Result()
	if err != nil {
		return nil, err
	}

	// 验证码判断
	if in.VerificationCode != result {
		return nil, errors.New("短信验证码错误,请重新输入")
	}

	// 密码一致性检查
	if in.Password != in.ConfirmPassword {
		return nil, errors.New("两次输入的密码不一致,请重新输入")
	}

	// 更新密码
	user.Password = utils.Md5(in.Password)
	if err := user.Editor(global.DB); err != nil {
		return nil, err
	}

	return &proto.ForgotPasswordResp{}, nil
}

// CompleteInformation 完善信息
//
// 参数:
//   - ctx context.Context: 上下文
//   - in *proto.CompleteInformationReq: 完善信息请求，包含用户ID、头像、昵称、性别、真实姓名、身份证、出生日期
//
// 返回值:
//   - *proto.CompleteInformationResp: 完善信息响应
//   - error: 错误信息
func (s *Server) CompleteInformation(ctx context.Context, in *proto.CompleteInformationReq) (*proto.CompleteInformationResp, error) {
	// 获取用户信息
	var user model.ShunUser
	if err := user.GetUserById(global.DB, int(in.UserId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// 如果用户上传了身份证照片，使用OCR自动识别信息
	if in.IdCardFrontUrl != "" {
		//上传身份证正面照片进行识别，不记录在数据库中
		ocrResult, err := utils.AliOCR(in.IdCardFrontUrl, "id-card-front")
		if err == nil {
			// 解析 OCR 识别结果
			parsedResult, parseErr := utils.ParseOCRResult(ocrResult, "id-card-front")
			//已填信息比对
			if parseErr == nil {
				// 验证 OCR 识别结果与用户填写信息是否一致
				if parsedResult.RealName != "" && in.RealName != "" && parsedResult.RealName != in.RealName {
					return nil, errors.New("身份证姓名与上传照片信息不一致")
				}
				if parsedResult.Birthday != "" && in.BirthdayTime != "" && parsedResult.Birthday != in.BirthdayTime {
					return nil, errors.New("出生日期与上传照片信息不一致")
				}
				if parsedResult.Gender != "" && in.Sex != "" && parsedResult.Gender != in.Sex {
					return nil, errors.New("性别与上传照片信息不一致")
				}

				// 使用 OCR 识别结果填充字段
				if parsedResult.RealName != "" {
					in.RealName = parsedResult.RealName
				}
				if parsedResult.IdCard != "" {
					in.IdCard = parsedResult.IdCard
				}
				if parsedResult.Birthday != "" {
					in.BirthdayTime = parsedResult.Birthday
				}
				if parsedResult.Gender != "" {
					in.Sex = parsedResult.Gender
				}
			}
		}
	}

	// 身份证验证
	if in.RealName != "" && in.IdCard != "" {
		isValid, err := utils.VerifyIdCard(in.RealName, in.IdCard)
		if err != nil {
			return nil, fmt.Errorf("身份证验证失败: %v", err)
		}
		if !isValid {
			return nil, errors.New("身份证信息不匹配")
		}
	}

	// 创建用户信息对象
	newUser := &model.ShunUser{
		Cover:        in.Cover,                                        // 头像
		Nickname:     in.Nickname,                                     // 昵称
		Sex:          in.Sex,                                          // 性别
		RealName:     utils.EnPwdCode([]byte(in.RealName)),            // 真实姓名
		IdCard:       utils.EnPwdCode([]byte(in.IdCard)),              // 身份证号
		BirthdayTime: utils.StringTransformationTime(in.BirthdayTime), // 出生日期
	}

	// 更新用户信息
	if err := newUser.Editor(global.DB); err != nil {
		return nil, err
	}

	// 删除缓存，保证数据一致性
	cacheKey := fmt.Sprintf("user:personal_center:%d", in.UserId)
	global.Rdb.Del(ctx, cacheKey)

	return &proto.CompleteInformationResp{}, nil
}

// StudentVerification 学生认证
//
// 参数:
//   - ctx context.Context: 上下文
//   - in *proto.StudentVerificationReq: 学生认证请求，包含用户ID、真实姓名、学号、学校名称、入学年份、学生证照片
//
// 返回值:
//   - *proto.StudentVerificationResp: 学生认证响应
//   - error: 错误信息
func (s *Server) StudentVerification(ctx context.Context, in *proto.StudentVerificationReq) (*proto.StudentVerificationResp, error) {
	// 获取用户信息
	var user model.ShunUser
	if err := user.GetUserById(global.DB, int(in.UserId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// 学生优惠过期时间
	var studentDiscountExpireTime time.Time

	// 如果用户上传了学生证照片，使用OCR自动识别信息
	if in.StudentIdPhoto != "" {
		// 这里假设 in.StudentIdPhoto 是学生证照片的URL
		// 实际使用时，需要根据前端传递的参数进行调整
		ocrResult, err := utils.AliOCR(in.StudentIdPhoto, "student-card")
		if err == nil {
			// 解析 OCR 识别结果
			parsedResult, parseErr := utils.ParseOCRResult(ocrResult, "student-card")
			if parseErr == nil {
				// 验证 OCR 识别结果与用户填写信息是否一致
				if parsedResult.SchoolName != "" && in.SchoolName != "" && parsedResult.SchoolName != in.SchoolName {
					return nil, errors.New("学校名称与上传照片信息不一致")
				}
				if parsedResult.StudentId != "" && in.StudentId != "" && parsedResult.StudentId != in.StudentId {
					return nil, errors.New("学号与上传照片信息不一致")
				}

				// 使用 OCR 识别结果填充字段
				if parsedResult.SchoolName != "" {
					in.SchoolName = parsedResult.SchoolName
				}
				if parsedResult.StudentId != "" {
					in.StudentId = parsedResult.StudentId
				}

				// 处理学生证到期时间
				if parsedResult.ExpireDate != "" {
					// 解析学生证到期时间
					studentCardExpireTime, parseErr := time.Parse("20060102", parsedResult.ExpireDate)
					if parseErr == nil {
						// 学生优惠过期时间 = 学生证到期时间 + 1年
						studentDiscountExpireTime = studentCardExpireTime.AddDate(1, 0, 0)
					}
				}
			}
		}
	}

	// 验证学生认证信息与用户实名认证信息一致
	if user.RealName != "" && user.IdCard != "" {
		if in.RealName != "" && in.RealName != user.RealName {
			return nil, errors.New("学生认证姓名与用户实名认证姓名不一致")
		}
	}

	// 学生认证
	if in.RealName != "" && in.StudentId != "" && in.SchoolName != "" {
		isValid, err := utils.StudentVerification(in.RealName, in.StudentId, in.SchoolName)
		if err != nil {
			return nil, fmt.Errorf("学生认证失败: %v", err)
		}
		if !isValid {
			return nil, errors.New("学生信息不匹配")
		}
	}

	// 创建用户信息对象
	newUser := &model.ShunUser{
		SchoolName:     in.SchoolName,                                     // 学校名称
		StudentId:      utils.EnPwdCode([]byte(in.StudentId)),             // 学号
		EnrollmentYear: utils.StringTransformationTime(in.EnrollmentYear), // 入学年份
		StudentIdPhoto: utils.EnPwdCode([]byte(in.StudentIdPhoto)),        // 学生证照片
	}

	// 设置学生优惠过期时间（从OCR识别结果计算）
	if !studentDiscountExpireTime.IsZero() {
		newUser.ExpirationTime = studentDiscountExpireTime
	}

	// 更新用户信息
	if err := newUser.Editor(global.DB); err != nil {
		return nil, err
	}

	// 删除缓存，保证数据一致性
	cacheKey := fmt.Sprintf("user:personal_center:%d", in.UserId)
	global.Rdb.Del(ctx, cacheKey)

	return &proto.StudentVerificationResp{}, nil
}

// PersonalCenter 个人中心
//
// 参数:
//   - ctx context.Context: 上下文
//   - in *proto.PersonalCenterReq: 个人中心请求，包含用户ID
//
// 返回值:
//   - *proto.PersonalCenterResp: 个人中心响应，包含用户信息
//   - error: 错误信息
func (s *Server) PersonalCenter(ctx context.Context, in *proto.PersonalCenterReq) (*proto.PersonalCenterResp, error) {
	// 生成缓存键
	cacheKey := fmt.Sprintf("user:personal_center:%d", in.UserId)

	// 尝试从缓存获取
	cachedData, err := global.Rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中，解析数据
		var resp proto.PersonalCenterResp
		if json.Unmarshal([]byte(cachedData), &resp) == nil {
			return &resp, nil
		}
	}

	// 缓存未命中，从数据库获取
	var user model.ShunUser
	if err := user.GetUserById(global.DB, int(in.UserId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// 构建响应
	resp := &proto.PersonalCenterResp{
		Phone:        utils.PhoneDesensitization(string(utils.DePwdCode(user.Phone))), // 手机号脱敏
		Cover:        user.Cover,                                                      // 头像
		Nickname:     user.Nickname,                                                   // 昵称
		Sex:          user.Sex,                                                        // 性别
		BirthdayTime: utils.TimeTransformationString(user.BirthdayTime),               // 出生日期
	}

	// 存入缓存，设置过期时间
	if data, err := json.Marshal(resp); err == nil {
		global.Rdb.Set(ctx, cacheKey, data, time.Minute*30)
	}

	return resp, nil
}

// Logout 账号注销
//
// 参数:
//   - _ context.Context: 上下文（未使用）
//   - in *proto.LogoutReq: 注销请求，包含用户ID
//
// 返回值:
//   - *proto.LogoutResp: 注销响应
//   - error: 错误信息
func (s *Server) Logout(_ context.Context, in *proto.LogoutReq) (*proto.LogoutResp, error) {
	// 获取用户信息
	var user model.ShunUser
	if err := user.GetUserById(global.DB, int(in.UserId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账号不存在")
		}
	}

	// 删除用户信息
	if err := user.Remove(global.DB); err != nil {
		return nil, err
	}

	// 将账号状态改为3注销
	user.Status = "3"
	if err := user.Editor(global.DB); err != nil {
		return nil, err
	}

	return &proto.LogoutResp{}, nil
}

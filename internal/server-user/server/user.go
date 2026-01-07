package server

import (
	"context"
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

type Server struct {
	proto.UnimplementedUserServer
}

// SendTextMessage 短信发送
func (s *Server) SendTextMessage(_ context.Context, in *proto.SendTextMessageReq) (*proto.SendTextMessageResp, error) {
	//1分钟发送次数Key
	minuteSendCountKey := fmt.Sprintf("%s:minute_send_text-message_count", in.Phone)
	count, _ := global.Rdb.Get(context.Background(), minuteSendCountKey).Int() //获取发送次数
	// >=1 发送拦截
	if count >= 1 {
		return nil, errors.New("一分钟内只能发送1次短信")
	}
	//1小时发送次数key
	hourSendCountKey := fmt.Sprintf("%s:hour_send_text-message_count", in.Phone)
	counts, _ := global.Rdb.Get(context.Background(), hourSendCountKey).Int()
	// >=5发送拦截
	if counts >= 5 {
		return nil, errors.New("一小时内只能发送5次短信")
	}
	//随机数生成验证码
	code := rand.Intn(9000) + 1000
	//互亿第三方短信平台验证
	utils.HuYi(in.Phone, strconv.Itoa(code))
	//设置发送短信验证码的键
	sendKey := fmt.Sprintf("%s:send_text-message", in.Phone)
	global.Rdb.Set(context.Background(), sendKey, code, time.Minute*5) //存入redis,5分钟过期时间
	//记录一分钟内的发送次数
	incr := global.Rdb.Incr(context.Background(), minuteSendCountKey)
	if incr.Val() == 1 {
		global.Rdb.Expire(context.Background(), minuteSendCountKey, time.Minute*1) //设置1分钟过期时间
	}
	//记录1小时内的发送次数
	hourIncr := global.Rdb.Incr(context.Background(), hourSendCountKey)
	if hourIncr.Val() == 1 {
		global.Rdb.Expire(context.Background(), hourSendCountKey, time.Hour*1) //设置1小时的过期时间
	}
	return &proto.SendTextMessageResp{}, nil
}

// Register 注册(自动登录)
func (s *Server) Register(_ context.Context, in *proto.RegisterReq) (*proto.RegisterResp, error) {
	//账号查询，防止重复账号注册
	var user model.ShunUser
	if err := user.GetUserByPhone(global.DB, in.Phone); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账号已存在")
		}
	}
	//从缓存中获取短信验证码
	sendKey := fmt.Sprintf("%s:send_text-message", in.Phone)
	result, err := global.Rdb.Get(context.Background(), sendKey).Result()
	if err != nil {
		return nil, err
	}
	//验证码判断
	if in.VerificationCode != result {
		return nil, errors.New("短信验证码错误,请重新输入")
	}
	user.Phone = in.Phone
	user.Password = in.Password
	user.Cover = "https://gd-hbimg.huaban.com/9980f53d02800ebad0472f0fe1a6eed9a1ba699ec27f-mvUjPb_fw658" // 设置默认头像
	if err := user.CreateUser(global.DB); err != nil {
		return nil, err
	}
	// 设置昵称，格式为：用户+注册时间+用户ID
	user.Nickname = "用户" + time.Now().Format("20060102150405") + strconv.FormatUint(user.Id, 10)
	//注册成功后自动登录并更新最后登录时间
	user.LastLoginTime = time.Now()
	if err := user.Editor(global.DB); err != nil {
		return nil, err
	}
	//最后删除缓存中的验证码，防止重复使用
	global.Rdb.Del(context.Background(), sendKey)
	return &proto.RegisterResp{Id: int64(user.Id)}, nil
}

// Login 登录
func (s *Server) Login(_ context.Context, in *proto.LoginReq) (*proto.LoginResp, error) {
	//账号查询判断
	var user model.ShunUser
	if err := user.GetUserByPhone(global.DB, in.Phone); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账号不存在")
		}
	}
	//判断登录方式
	if in.Password == "" && in.VerificationCode != "" { //短信验证码登录
		//从缓存中获取短信验证码
		sendKey := fmt.Sprintf("%s:send_text-message", in.Phone)
		result, err := global.Rdb.Get(context.Background(), sendKey).Result()
		if err != nil {
			return nil, err
		}
		//验证码判断
		if in.VerificationCode != result {
			return nil, errors.New("短信验证码错误,请重新输入")
		}
		//最后删除缓存中的验证码，防止重复使用
		global.Rdb.Del(context.Background(), sendKey)
	} else { //密码登录
		banKey := fmt.Sprintf("%s:login_ban", in.Phone)                  //设置错误次数键
		count, err := global.Rdb.Get(context.Background(), banKey).Int() //获取
		if err != nil {
			return nil, err
		}
		//错误次数>=3登录拦截
		if count >= 3 {
			return nil, errors.New("登录次数过多,请2小时后在登录")
		}
		if user.Password != in.Password { //密码错误
			incr := global.Rdb.Incr(context.Background(), banKey) //记录错误次数
			if incr.Val() == 1 {
				global.Rdb.Expire(context.Background(), banKey, time.Hour*2) //设置过期时间2小时
			}
			return nil, errors.New("密码错误,请重新输入")
		}
		//登录成功,删除错误次数键
		global.Rdb.Del(context.Background(), banKey)
	}
	//更新最后登录时间
	user.LastLoginTime = time.Now()
	if err := user.Editor(global.DB); err != nil {
		return nil, err
	}
	return &proto.LoginResp{Id: int64(user.Id)}, nil
}

// ForgotPassword 忘记密码
func (s *Server) ForgotPassword(_ context.Context, in *proto.ForgotPasswordReq) (*proto.ForgotPasswordResp, error) {
	//账号查询判断
	var user model.ShunUser
	if err := user.GetUserByPhone(global.DB, in.Phone); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账号不存在")
		}
	}
	//从缓存中获取短信验证码
	sendKey := fmt.Sprintf("%s:send_text-message", in.Phone)
	result, err := global.Rdb.Get(context.Background(), sendKey).Result()
	if err != nil {
		return nil, err
	}
	//验证码判断
	if in.VerificationCode != result {
		return nil, errors.New("短信验证码错误,请重新输入")
	}
	user.Password = in.Password
	if in.Password != in.ConfirmPassword {
		return nil, errors.New("两次输入的密码不一致,请重新输入")
	}
	if err := user.Editor(global.DB); err != nil {
		return nil, err
	}
	return &proto.ForgotPasswordResp{}, nil
}

// CompleteInformation 完善信息
func (s *Server) CompleteInformation(_ context.Context, in *proto.CompleteInformationReq) (*proto.CompleteInformationResp, error) {
	var user model.ShunUser
	if err := user.GetUserById(global.DB, int(in.UserId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	newUser := &model.ShunUser{
		Cover:        in.Cover,
		Nickname:     in.Nickname,
		Sex:          in.Sex,
		RealName:     in.RealName,
		IdCard:       in.IdCard,
		BirthdayTime: utils.StringTransformationTime(in.BirthdayTime),
	}
	if err := newUser.Editor(global.DB); err != nil {
		return nil, err
	}
	return &proto.CompleteInformationResp{}, nil
}

// StudentVerification 学生认证
func (s *Server) StudentVerification(_ context.Context, in *proto.StudentVerificationReq) (*proto.StudentVerificationResp, error) {
	var user model.ShunUser
	if err := user.GetUserById(global.DB, int(in.UserId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	newUser := &model.ShunUser{
		SchoolName:     in.SchoolName,
		StudentId:      in.StudentId,
		EnrollmentYear: utils.StringTransformationTime(in.EnrollmentYear),
		StudentIdPhoto: in.StudentIdPhoto,
	}
	if err := newUser.Editor(global.DB); err != nil {
		return nil, err
	}
	return &proto.StudentVerificationResp{}, nil
}

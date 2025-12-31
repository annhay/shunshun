package server

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"shunshun/internal/pkg/global"
	"shunshun/internal/pkg/utils"
	"shunshun/internal/proto"
	"strconv"
	"time"
)

type Server struct {
	proto.UnimplementedUserServer
}

// SendTextMessage 短信发送
func (s *Server) SendTextMessage(_ context.Context, in *proto.SendTextMessageReq) (*proto.SendTextMessageResp, error) {
	//发送次数Key
	sendCountKey := fmt.Sprintf("%ssend_text-message_count", in.Phone)
	count, _ := global.Rdb.Get(context.Background(), sendCountKey).Int() //获取发送次数
	// >=1 发送拦截
	if count >= 1 {
		return nil, errors.New("一分钟内只能发送一次短信")
	}
	//随机数生成验证码
	code := rand.Intn(9000) + 1000
	//互亿第三方短信平台验证
	_, err := utils.HuYi(in.Phone, strconv.Itoa(code))
	if err != nil {
		return nil, err
	}
	//发送短信验证码Key
	sendKey := fmt.Sprintf("%ssend_text-message", in.Phone)
	global.Rdb.Set(context.Background(), sendKey, code, time.Minute*5) //存入redis,5分钟过期时间
	//发送次数+1
	incr := global.Rdb.Incr(context.Background(), sendCountKey)
	if incr.Val() >= 1 {
		//设置1分钟过期时间
		global.Rdb.Expire(context.Background(), sendCountKey, time.Minute*1)
	}
	return &proto.SendTextMessageResp{}, nil
}

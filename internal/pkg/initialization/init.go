package initialization

import (
	"shunshun/internal/pkg/global"
	"shunshun/internal/pkg/utils"

	"go.uber.org/zap"
)

// GatewayInit 网关层初始化
//
// 初始化内容:
//  1. 用户服务gRPC客户端
//  2. 司机服务gRPC客户端
//  3. 订单服务gRPC客户端
func GatewayInit() {
	UserGrpcClient()   // 初始化用户服务gRPC客户端
	DriverGrpcClient() // 初始化司机服务gRPC客户端
	OrderGrpcClient()  // 初始化订单服务gRPC客户端
}

// ServerInit 服务初始化
//
// 初始化内容:
//  1. 配置文件加载
//  2. 日志系统初始化
//  3. MySQL数据库连接
//  4. Redis缓存连接
//  5. RabbitMQ消息队列
func ServerInit() {
	ViperInit() // 初始化Viper配置
	// 初始化Zap日志
	global.Logger = InitLogger()
	defer global.Logger.Sync() // 延迟关闭日志，确保日志写入
	MysqlInit()                // 初始化MySQL数据库
	RedisInit()                // 初始化Redis缓存

	// 初始化RabbitMQ
	rmqConfig := utils.RabbitMQConfig{
		Host:       global.AppConf.RabbitMQ.Host,
		Port:       global.AppConf.RabbitMQ.Port,
		User:       global.AppConf.RabbitMQ.User,
		Password:   global.AppConf.RabbitMQ.Password,
		VHost:      global.AppConf.RabbitMQ.VHost,
		Exchange:   global.AppConf.RabbitMQ.Exchange,
		Queue:      global.AppConf.RabbitMQ.Queue,
		RoutingKey: global.AppConf.RabbitMQ.RoutingKey,
	}

	err := utils.InitRabbitMQ(rmqConfig)
	if err != nil {
		global.Logger.Error("Failed to initialize RabbitMQ", zap.Error(err))
		// 这里不panic，因为RabbitMQ可能暂时不可用，但系统其他功能仍可运行
	}
}

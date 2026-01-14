package initialization

import (
	"shunshun/internal/pkg/global"
)

// GatewayInit 网关层初始化
func GatewayInit() {
	UserGrpcClient()
	DriverGrpcClient()
}

// ServerInit 服务初始化
func ServerInit() {
	ViperInit()
	// 初始化Zap日志
	global.Logger = InitLogger()
	defer global.Logger.Sync()
	MysqlInit()
	RedisInit()
}

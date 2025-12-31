package initialization

// GatewayInit 网关层初始化
func GatewayInit() {
	UserGrpcClient()
}

// ServerInit 服务初始化
func ServerInit() {
	ViperInit()
	MysqlInit()
	RedisInit()
}

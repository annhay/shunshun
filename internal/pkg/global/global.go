package global

import (
	"context"
	"shunshun/internal/pkg/configs"
	"shunshun/internal/proto"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	AppConf      *configs.AppConfig
	DB           *gorm.DB
	Rdb          *redis.Client
	Logger       *zap.Logger
	UserClient   proto.UserClient
	DriverClient proto.DriverClient
	OrderClient  proto.OrderClient
	Ctx          = context.Background()
)

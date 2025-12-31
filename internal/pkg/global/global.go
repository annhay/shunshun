package global

import (
	"shunshun/internal/pkg/configs"
	"shunshun/internal/proto"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	AppConf    *configs.AppConfig
	DB         *gorm.DB
	Rdb        *redis.Client
	UserClient proto.UserClient
)

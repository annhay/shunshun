package initialization

import (
	"context"
	"fmt"
	"shunshun/internal/pkg/global"

	"github.com/redis/go-redis/v9"
)

func RedisInit() {
	var data = global.AppConf.Redis
	addr := fmt.Sprintf("%s:%d", data.Host, data.Port)
	global.Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: data.Password, // no password set
		DB:       data.DB,       // use default DB
	})
	err := global.Rdb.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
	global.Logger.Info("redis init success")
}

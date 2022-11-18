package initialize

import (
	"fmt"
	"github.com/go-redis/redis/v9"
	"web-api/good-web/global"
)

func InitRedis() {
	global.Rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host,
			global.ServerConfig.RedisInfo.Port),
	})
	if global.Rdb == nil {
		panic("can not connect redis")
	}
}

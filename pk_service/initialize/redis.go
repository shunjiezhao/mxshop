package initialize

import (
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"server/pk_service/global"
)

func InitRedis() {
	global.RedisClient = goredislib.NewClient(&goredislib.Options{
		Addr: global.Settings.RedisInfo.Addr,
	})
	pool := goredis.NewPool(global.RedisClient) // or, pool := redigo.NewPool(...)
	global.RedisPool = redsync.New(pool)
}

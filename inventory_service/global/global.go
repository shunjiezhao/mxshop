package global

import (
	"crypto/md5"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"server/inventory_service/config"
	"server/inventory_service/utils/queue"
	"server/shared/etcd/register"
)

var (
	DB             *gorm.DB
	HashMethodName = "pbkdf2-sha512"
	SaltLen        = 9
	Iterations     = 99
	KeyLen         = 32
	Options        = &password.Options{SaltLen, Iterations, KeyLen, md5.New}
	Settings       = &config.ServiceConfig{}
)

// etcd
var (
	ServiceRegister *register.ServiceRegister
)

// redis
var (
	RedisPool *redsync.Redsync
	Rdb       *redis.Client
)

var (
	StockRebackPublisher  *queue.Publisher
	StockRebackSubscriber *queue.Subscriber
)

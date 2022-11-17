package global

import (
	"crypto/md5"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"server/shared/etcd"
	"server/shopcart_service/config"
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
	ServiceRegister *etcd.ServiceRegister
)

// redis
var (
	RedisPool *redsync.Redsync
)

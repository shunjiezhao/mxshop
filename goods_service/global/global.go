package global

import (
	"crypto/md5"
	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/gorm"
	"server/goods_service/config"
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

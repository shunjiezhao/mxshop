package global

import (
	"crypto/md5"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	proto "server/goods_service/api/gen/v1/goods"
	proto3 "server/inventory_service/proto/gen/v1/inventory"
	"server/shared/etcd/discovery"
	"server/shared/etcd/register"
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
	ServiceRegister *register.ServiceRegister
)

// redis
var (
	RedisPool *redsync.Redsync
)

var (
	GoodSrv         proto.GoodsClient
	InventorySrv    proto3.InventoryClient
	ServerDiscovery *discovery.ServiceDiscovery
)

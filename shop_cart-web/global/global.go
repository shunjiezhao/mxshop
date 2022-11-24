package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v9"
	proto2 "web-api/good-web/proto"
	"web-api/shared/etcd/discovery"
	"web-api/shop_cart-web/config"
	"web-api/shop_cart-web/proto"
	"web-api/shop_cart-web/utils/token"
)

var (
	ServerConfig *config.ServiceConfig = &config.ServiceConfig{}
	Trans        ut.Translator
	OrderClient  proto.OrderClient
	GoodsClient  proto2.GoodsClient
)

//token
var (
	JWTTokenVerifier = &token.JWTTokenVerifier{}
)

// redis
var (
	Rdb *redis.Client
)

// ectd
var (
	ServerDiscovery *discovery.ServiceDiscovery
)

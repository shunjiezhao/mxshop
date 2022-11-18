package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v9"
	"web-api/good-web/config"
	goodpb "web-api/good-web/proto"
	"web-api/good-web/utils/token"
	"web-api/share/etcd/discovery"
)

var (
	ServerConfig       *config.ServiceConfig = &config.ServiceConfig{}
	Trans              ut.Translator
	GoodsServiceClient goodpb.GoodsClient
)

//token
var (
	JWTTokenVerifier = &token.JWTTokenVerifier{}
	JwtTokenGen      = &token.JWTokenGen{}
)

// redis
var (
	Rdb *redis.Client
)

// ectd
var (
	ServerDiscovery *discovery.ServiceDiscovery
)

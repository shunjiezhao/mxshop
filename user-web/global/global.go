package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v9"
	"web-api/user-web/config"
	"web-api/user-web/etcd/discovery"
	userpb "web-api/user-web/proto"
	"web-api/user-web/utils/token"
)

var (
	ServerConfig      *config.ServiceConfig = &config.ServiceConfig{}
	Trans             ut.Translator
	UserServiceClient userpb.UserServiceClient
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

package global

import (
	ut "github.com/go-playground/universal-translator"
	"web-api/user-web/config"
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

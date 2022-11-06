package global

import (
	ut "github.com/go-playground/universal-translator"
	"web-api/user-web/config"
)

var (
	ServerConfig *config.ServiceConfig = &config.ServiceConfig{}
	Trans        ut.Translator
)

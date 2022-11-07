package config

import (
	"time"
)

// 可以使用 mapstructure 来制定 配置文件的相应字段
type UserSettingS struct {
	Host string
	Port int
}
type JWTSettingS struct {
	Secret         string
	Issuer         string
	Expire         time.Duration
	PublicKeyPath  string
	PrivateKeyPath string
}

type ServiceConfig struct {
	Name        string
	GinAddr     string       // gin 运行端口
	UserSrvInfo UserSettingS `mapstructure:"UserSrv"`
	JwtInfo     JWTSettingS  `mapstructure:"JwtConfig"`
}

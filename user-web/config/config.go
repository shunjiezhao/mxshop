package config

import (
	"time"
)

type Addr struct {
	Host string `mapstructure:"Host" json:"host"`
	Port int    `mapstructure:"Port" json:"port"`
}

// 可以使用 mapstructure 来制定 配置文件的相应字段
type UserSettingS struct {
	Host string `mapstructure:"Host" json:"host"`
	Port int    `mapstructure:"Port" json:"port"`
}
type JWTSettingS struct {
	Secret         string
	Issuer         string
	ExpireMin      time.Duration
	PublicKeyPath  string `yaml:"PublicKeyPath"`
	PrivateKeyPath string `yaml:"PrivateKeyPath"`
}
type RedisSettings struct {
	Host      string `mapstructure:"Host" json:"host"`
	Port      int    `mapstructure:"Port" json:"port"`
	ExpireMin time.Duration
}
type EtcdSettings struct {
	EndPoints []string `yaml:"EndPoints"`
	Prefix    string   `mapstructure:"Prefix"`
}
type ServiceConfig struct {
	Name        string
	GinAddr     string        // gin 运行端口
	SrvName     string        `mapstructure:"SrvName"` // 这是 grpc 服务的名字 schema + srvName
	Debug       bool          `mapstructure:"Debug"`
	UserSrvInfo UserSettingS  `mapstructure:"UserSrv"`
	RedisInfo   RedisSettings `mapstructure:"RedisConfig"`
	JwtInfo     JWTSettingS   `mapstructure:"JwtConfig"`
	EtcdInfo    EtcdSettings  `mapstructure:"EtcdConfig"`
}

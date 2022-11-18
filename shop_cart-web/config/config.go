package config

import (
	"time"
)

type Addr struct {
	Host string `mapstructure:"Host" json:"host"`
	Port int    `mapstructure:"Port" json:"port"`
}

type JWTSettingS struct {
	Secret         string
	Issuer         string
	ExpireMin      time.Duration
	PublicKeyPath  string
	PrivateKeyPath string
}
type RedisSettings struct {
	Host      string `mapstructure:"Host" json:"host"`
	Port      int    `mapstructure:"Port" json:"port"`
	ExpireMin time.Duration
}
type EtcdSettings struct {
	EndPoints        []string `yaml:"EndPoints"`
	Prefix           string   `mapstructure:"Prefix"`
	GoodsSrvName     string   `yaml:"GoodsSrvName"`
	InventorySrvName string   `yaml:"InventorySrvName"`
	UserSrvName      string   `yaml:"UserSrvName"`
	CartSrvName      string   `yaml:"CartSrvName"`
}
type ServiceConfig struct {
	Name      string
	GinAddr   string        // gin 运行端口
	SrvName   string        `mapstructure:"SrvName"` // 这是 grpc 服务的名字 schema + srvName
	Debug     bool          `mapstructure:"Debug"`
	RedisInfo RedisSettings `mapstructure:"RedisConfig"`
	JwtInfo   JWTSettingS   `mapstructure:"JwtConfig"`
	EtcdInfo  EtcdSettings  `mapstructure:"EtcdConfig"`
}

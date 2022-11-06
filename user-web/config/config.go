package config

// 可以使用 mapstructure 来制定 配置文件的相应字段
type UserSrvConfig struct {
	Host string
	Port int
}
type ServiceConfig struct {
	Name        string
	GinAddr     string        // gin 运行端口
	UserSrvInfo UserSrvConfig `mapstructure:"UserSrv"`
}

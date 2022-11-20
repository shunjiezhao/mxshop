package config

type DBSettingS struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	User     string `mapstructure:"user" json:"user"`
	Dbname   string `mapstructure:"dbname" json:"dbname"`
	Password string `mapstructure:"password" json:"password"`
}

type Addr struct {
	Host string `mapstructure:"Host" json:"host"`
	Port int    `mapstructure:"Port" json:"port"`
}

type EtcdSettings struct {
	EndPoints []string `yaml:"EndPoints"`
	Prefix    string   `mapstructure:"Prefix"`
	LeaseSec  int64    `mapstructure:"LeaseSec"`

	GoodSrvName      string `yaml:"GoodSrvName"`
	InventorySrvName string `yaml:"InventorySrvName"`
}
type RedisSettings struct {
	Addr string `yaml:"Addr"`
}

type ServiceConfig struct {
	// 服务器监听的 grpc 端口
	IP      string `mapstructure:"IP"`
	Port    int    `mapstructure:"Port"`
	SrvName string `mapstructure:"SrvName"`

	DBConfig DBSettingS   `json:"db_config" mapstructure:"DBConfig"`
	EtcdInfo EtcdSettings `mapstructure:"EtcdConfig"`
	RedisInfo RedisSettings `yaml:"RedisConfig"`
}

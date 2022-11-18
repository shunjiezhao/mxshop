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
	Prefix    string   `yaml:"Prefix"`
	LeaseSec  int64    `yaml:"LeaseSec"`
}

type ServiceConfig struct {
	// 服务器监听的 grpc 端口
	IP      string `yaml:"IP"`
	Port    int    `yaml:"Port"`
	SrvName string `yaml:"SrvName"`

	DBConfig DBSettingS   `json:"db_config" mapstructure:"DBConfig"`
	EtcdInfo EtcdSettings `mapstructure:"EtcdConfig"`
}

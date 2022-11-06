package initialize

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"web-api/user-web/global"
)

var vp *viper.Viper

func getEnvBool(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
func InitConfig() {
	debug := getEnvBool("MXSHOP_DEBUG")
	configFileName := "user-web/config/config-pro.yaml"
	if debug {
		configFileName = "user-web/config/config-debug.yaml"
	}
	err := readFile(configFileName, "yaml")
	if err != nil {
		panic(err)
	}
	readSection("UserWeb", global.ServerConfig)
	zap.L().Info("得到配置", zap.String("配置文件名", configFileName),
		zap.Any("配置信息", global.ServerConfig))

	fmt.Println("%v\n", global.ServerConfig)
}

func readFile(path, fileType string) error {
	vp = viper.New()
	//vp.AddConfigPath("")              // path
	vp.SetConfigFile(path)     // filename
	vp.SetConfigType(fileType) // .type
	err := vp.ReadInConfig()
	return err
}

func readSection(k string, v interface{}) error {
	err := vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}

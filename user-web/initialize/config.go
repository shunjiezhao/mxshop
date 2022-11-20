package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"web-api/user-web/global"
	"web-api/user-web/utils"
)

var vp *viper.Viper

var (
	prevHash string
)

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
	prevHash = getFileHash(path)

	vp.WatchConfig()
	// 如果改变的话，需要动态加载一变
	vp.OnConfigChange(func(in fsnotify.Event) {
		zap.L().Info("config file is change", zap.Any("in", in))
		if in.Op == fsnotify.Write {
			if nowHash := getFileHash(path); nowHash != "" && nowHash == prevHash {
				zap.L().Info("don't change")
				return
			}
			zap.L().Info("config file is change")
			InitValidator("zh")
			InitConnect()
			InitJwtVerifier()
			InitRedis()
		}
	})
	return err
}

func readSection(k string, v interface{}) error {
	err := vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}

func getFileHash(path string) string {
	open, err := os.Open(path)
	if err != nil {
		zap.L().Error("can not open file", zap.String("path", path))
		return ""
	}
	b, _ := ioutil.ReadAll(open)
	return utils.Hash(b)
}

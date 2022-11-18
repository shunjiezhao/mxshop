package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"web-api/good-web/global"
	"web-api/good-web/initialize"
)

func main() {
	// 初始化工作
	initialize.InitLogger()
	// 初始化配置文件
	initialize.InitConfig()
	initialize.InitValidator("zh")
	initialize.InitConnect()
	initialize.InitJwtVerifier()
	initialize.InitRedis()
	routers := initialize.Routers()

	log.Fatal(routers.Run(global.ServerConfig.GinAddr))

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// 注销服务
}

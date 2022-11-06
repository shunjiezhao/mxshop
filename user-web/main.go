package main

import (
	"log"
	"web-api/user-web/global"
	"web-api/user-web/initialize"
)

func main() {
	// 初始化工作
	initialize.InitLogger()
	// 初始化配置文件
	initialize.InitConfig()
	initialize.InitValidator("zh")

	routers := initialize.Routers()

	log.Fatal(routers.Run(global.ServerConfig.GinAddr))
}

package main

import (
	"context"
	"log"
	"time"
	"web-api/user-web/global"
	"web-api/user-web/initialize"
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

	global.Rdb.Set(context.Background(), "123", 1, 1*time.Minute)

	routers := initialize.Routers()
	log.Fatal(routers.Run(global.ServerConfig.GinAddr))
}

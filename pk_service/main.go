package main

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"server/pk_service/global"
	"server/pk_service/handler"
	"server/pk_service/initialize"
	proto "server/pk_service/proto/gen/v1/pk"
	"syscall"
)

func main() {
	initialize.InitConfig() // 初始化配置;
	initialize.InitDB()
	initialize.InitRedis()
	initialize.InitConnect()
	initialize.InitQueue()

	address := fmt.Sprintf("%s:%d", global.Settings.IP, global.Settings.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("can not create tcp listener: %v", err)
	}

	logger, err := NewZapLogger()
	if err != nil {
		panic(err)
	}

	svc := grpc.NewServer()
	proto.RegisterPKServer(svc, handler.NewService(&handler.Config{
		DB:            global.DB,
		Logger:        logger,
		UserPublisher: global.UserWaitQueue,
		//TODO: 实现剩下两个接口
	}))

	logger.Info("grpc service run start", zap.String("name", "user"), zap.String("address", address))

	var etcClose io.Closer
	etcClose = initialize.InitEtcd(logger)
	go svc.Serve(lis)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// 注销服务
	err = etcClose.Close()
	if err != nil {
		zap.L().Error("注销失败")
		return
	}
	zap.L().Info("注销成功")

}
func NewZapLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = ""
	return cfg.Build()
}

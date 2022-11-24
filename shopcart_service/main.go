package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"server/shopcart_service/global"
	"server/shopcart_service/handler"
	"server/shopcart_service/initialize"
	proto "server/shopcart_service/proto/gen/v1/cart"
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
	service := handler.New(handler.Config{
		DB:     global.DB,
		Logger: logger,
	})
	proto.RegisterOrderServer(svc, service)
	//
	go service.Watch(context.Background())

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

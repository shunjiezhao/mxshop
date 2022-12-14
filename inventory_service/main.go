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
	"server/inventory_service/global"
	"server/inventory_service/handler"
	"server/inventory_service/initialize"
	proto "server/inventory_service/proto/gen/v1/inventory"
	"syscall"
)

func main() {
	initialize.InitConfig() // 初始化配置;
	initialize.InitDB()
	initialize.InitRedis()
	initialize.InitQueue()

	address := fmt.Sprintf("%s:%d", global.Settings.IP, global.Settings.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("can not create tcp listener: %v", err)
	}
	svc := grpc.NewServer()

	logger, err := NewZapLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("grpc service run start", zap.String("name", "user"), zap.String("address", address))

	service := handler.NewService(&handler.InventorySrvConfig{
		Logger:     logger,
		DB:         global.DB,
		Publisher:  global.StockRebackPublisher,
		Subscriber: global.StockRebackSubscriber,
	})
	proto.RegisterInventoryServer(svc, service)
	go service.WatchStockReback(context.Background())

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

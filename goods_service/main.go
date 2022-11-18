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
	proto "server/goods_service/api/gen/v1/goods"
	"server/goods_service/global"
	"server/goods_service/handler"
	"server/goods_service/initialize"
	"syscall"
)

func main() {
	initialize.InitConfig() // 初始化配置
	initialize.InitDB()

	address := fmt.Sprintf("%s:%d", global.Settings.IP, global.Settings.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("can not create tcp listener: %v", err)
	}
	svc := grpc.NewServer()
	proto.RegisterGoodsServer(svc, handler.New(global.DB))

	logger, err := NewZapLogger()
	if err != nil {
		panic(err)
	}
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

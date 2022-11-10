package main

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"server/good_service/global"
	"server/good_service/initialize"
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

	logger, err := NewZapLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("grpc service run start", zap.String("name", "user"), zap.String("address", address))

	go initialize.InitEtcd(logger)
	log.Fatalln(svc.Serve(lis))
}
func NewZapLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = ""
	return cfg.Build()
}

package main

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	proto "server/good_service/api/gen/v1/goods"
	"server/good_service/global"
	"server/good_service/handler"
	"server/good_service/initialize"
	"server/good_service/model"
	"time"
)

func main() {
	initialize.InitConfig() // 初始化配置
	initialize.InitDB()
	err := global.DB.AutoMigrate(
		&model.Category{},
		&model.Brands{},
		&model.GoodsCategoryBrand{},
		&model.Banner{},
		&model.Goods{},
	)
	if err != nil {
		panic(err)
	}
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if global.DB == nil {
					log.Println("db is nil")
				}
			}

		}
	}()

	address := fmt.Sprintf("%s:%d", global.Settings.IP, global.Settings.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("can not create tcp listener: %v", err)
	}
	svc := grpc.NewServer()
	proto.RegisterGoodsServer(svc, &handler.GoodsServer{})

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

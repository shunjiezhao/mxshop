package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	userpb "server/user_service/api/gen/v1/user"
	"server/user_service/dao"
	"server/user_service/global"
	"server/user_service/handler"
	"server/user_service/initialize"
	"server/user_service/model"
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

	global.DB.Migrator().DropTable(&model.User{})
	global.DB.AutoMigrate(&model.User{})

	service := handler.UserService{
		Dao: dao.New(global.DB),
	}

	for i := 1; i < 10; i++ {
		service.CreateUser(context.Background(), &userpb.CreateUserRequest{
			Mobile:   fmt.Sprintf("1334744689%d", i),
			PassWord: "123456",
			Nickname: fmt.Sprintf("name%d", i),
		})
	}
	userpb.RegisterUserServiceServer(svc, &handler.UserService{
		Dao: dao.New(global.DB),
	})
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

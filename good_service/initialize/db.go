package initialize

import (
	"fmt"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	proto "server/good_service/api/gen/v1/goods"
	"server/good_service/global"
	"server/good_service/handler"
)

func InitDB() {
	g := &gorm.Config{}
	var err error

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		global.Settings.DBConfig.Host,
		global.Settings.DBConfig.User,
		global.Settings.DBConfig.Dbname,
		global.Settings.DBConfig.Password,
		global.Settings.DBConfig.Port)

	global.DB, err = gorm.Open(postgres.Open(dsn), g)
	server := grpc.NewServer()
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})

	handlerErr(err)
	sqlDB, err := global.DB.DB()
	handlerErr(err)

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
}

// 初始化只要出错都 panic
func handlerErr(err error) {
	if err != nil {
		panic(err)
	}
}

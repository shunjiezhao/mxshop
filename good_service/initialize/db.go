package initialize

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"server/good_service/global"
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

package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"server/shared/etcd"
	"server/user_service/global"
)

func InitEtcd(logger *zap.Logger) {
	address := fmt.Sprintf("%s:%d", global.Settings.IP, global.Settings.Port)
	ser, err := etcd.NewServiceRegister(global.Settings.EtcdInfo.EndPoints, global.Settings.SrvName, address, global.Settings.EtcdInfo.LeaseSec, logger)

	go ser.Watch()
	if err != nil {
		logger.Error("Init Fail", zap.Error(err))
	}
	logger.Info("【Etcd】: 初始化成功")
}

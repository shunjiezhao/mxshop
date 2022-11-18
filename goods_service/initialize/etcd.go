package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"server/goods_service/global"
	"server/shared/etcd/register"
)

func InitEtcd(logger *zap.Logger) io.Closer {
	address := fmt.Sprintf("%s:%d", global.Settings.IP, global.Settings.Port)
	ser, err := register.NewServiceRegister(global.Settings.EtcdInfo.EndPoints, logger)

	//申请租约设置时间keepalive
	if err := ser.Register(global.Settings.EtcdInfo.Prefix, address, global.Settings.EtcdInfo.LeaseSec); err != nil {
		zap.L().Fatal("【Etcd】: 无法注册", zap.Error(err))
	}

	zap.L().Info("服务注册成功")
	go ser.Watch()
	if err != nil {
		logger.Error("Init Fail", zap.Error(err))
	}
	logger.Info("【Etcd】: 初始化成功")

	return ser

}

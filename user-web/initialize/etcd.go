package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"web-api/user-web/etcd/discovery"
	"web-api/user-web/global"
)

func InitEtcd() {
	var endpoints []string
	endpoints = append(endpoints, fmt.Sprintf("%s:%d", global.ServerConfig.EtcdInfo.Host, global.ServerConfig.EtcdInfo.Port))
	global.ServerDiscovery = discovery.NewServiceDiscovery(endpoints)
	err := global.ServerDiscovery.WatchService(global.ServerConfig.EtcdInfo.Prefix)
	if err != nil {
		zap.L().Error("Init Fail", zap.Error(err))
	}
}

package initialize

import (
	"google.golang.org/grpc/resolver"
	"web-api/good-web/global"
	"web-api/share/etcd/discovery"
)

func InitEtcd() {
	global.ServerDiscovery = discovery.NewServiceDiscovery(global.ServerConfig.EtcdInfo.EndPoints)
	resolver.Register(global.ServerDiscovery)
}
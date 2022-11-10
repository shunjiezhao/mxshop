package initialize

import (
	"google.golang.org/grpc/resolver"
	"web-api/user-web/etcd/discovery"
	"web-api/user-web/global"
)

func InitEtcd() {
	global.ServerDiscovery = discovery.NewServiceDiscovery(global.ServerConfig.EtcdInfo.EndPoints)
	resolver.Register(global.ServerDiscovery)
}

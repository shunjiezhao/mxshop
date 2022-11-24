package initialize

import (
	"google.golang.org/grpc/resolver"
	"web-api/shared/etcd/discovery"
	"web-api/shop_cart-web/global"
)

func InitEtcd() {
	global.ServerDiscovery = discovery.NewServiceDiscovery(global.ServerConfig.EtcdInfo.EndPoints)
	resolver.Register(global.ServerDiscovery)
}

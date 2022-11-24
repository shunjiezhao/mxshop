package initialize

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"web-api/good-web/global"
	"web-api/good-web/proto"
	"web-api/shared/etcd/discovery"
)

// 连接grpc
func InitConnect() {
	global.ServerDiscovery = discovery.NewServiceDiscovery(global.ServerConfig.EtcdInfo.EndPoints)
	resolver.Register(global.ServerDiscovery)
	conn, err := grpc.Dial(global.ServerDiscovery.Scheme()+"://zsj.com/"+global.ServerConfig.SrvName,
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("[GetGoodsList] 连接 【商品服务失败】", zap.Error(err))
		return
	}
	global.GoodsServiceClient = proto.NewGoodsClient(conn)
}

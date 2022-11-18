package initialize

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	proto2 "web-api/good-web/proto"
	"web-api/share/etcd/discovery"
	"web-api/shop_cart-web/global"
	"web-api/shop_cart-web/proto"
)

// 连接grpc
func InitConnect() {
	global.ServerDiscovery = discovery.NewServiceDiscovery(global.ServerConfig.EtcdInfo.EndPoints)
	resolver.Register(global.ServerDiscovery)
	conn, err := grpc.Dial(global.ServerDiscovery.Scheme()+"://zsj.com/"+global.ServerConfig.EtcdInfo.GoodsSrvName,
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("连接 【商品服务失败】", zap.Error(err))
		return
	}
	global.GoodsClient = proto2.NewGoodsClient(conn)

	conn, err = grpc.Dial(global.ServerDiscovery.Scheme()+"://zsj.com/"+global.ServerConfig.EtcdInfo.CartSrvName,
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("连接 【订单服务失败】", zap.Error(err))
		return
	}

	global.OrderClient = proto.NewOrderClient(conn)
}

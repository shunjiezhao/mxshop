package initialize

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	proto "server/goods_service/api/gen/v1/goods"
	proto2 "server/inventory_service/proto/gen/v1/inventory"
	"server/shared/etcd/discovery"
	"server/shopcart_service/global"
)

// 连接grpc
func InitConnect() {
	var err error
	global.ServerDiscovery = discovery.NewServiceDiscovery(global.Settings.EtcdInfo.EndPoints)
	resolver.Register(global.ServerDiscovery)
	conn, err := grpc.Dial(global.ServerDiscovery.Scheme()+"://zsj.com/"+global.Settings.EtcdInfo.GoodSrvName, grpc.WithInsecure())
	if err != nil {
		zap.L().Fatal("can not connect the goods service", zap.Error(err))
	}
	global.GoodSrv = proto.NewGoodsClient(conn)

	conn, err = grpc.Dial(global.ServerDiscovery.Scheme()+"://zsj.com/"+global.Settings.EtcdInfo.InventorySrvName, grpc.WithInsecure())
	if err != nil {
		zap.L().Fatal("can not connect the inventory service", zap.Error(err))
	}
	global.InventorySrv = proto2.NewInventoryClient(conn)
}

package initialize

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"web-api/user-web/etcd/discovery"
	"web-api/user-web/global"
	userpb "web-api/user-web/proto"
	"web-api/user-web/utils/divide"
)

// 连接grpc
func InitConnect() {
	global.ServerDiscovery = discovery.NewServiceDiscovery(global.ServerConfig.EtcdInfo.EndPoints)
	resolver.Register(global.ServerDiscovery)
	conn, err := grpc.Dial(global.ServerDiscovery.Scheme()+"://zsj.com/"+global.ServerConfig.SrvName,
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("[GetUserList] 连接 【用户服务失败】", zap.Error(err))
		return
	}
	global.UserServiceClient = userpb.NewUserServiceClient(conn)

	//TODO: 地址配置化
	conn, err = grpc.Dial(global.ServerDiscovery.Scheme()+"://zsj.com/"+"pk_srv",
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("[GetUserList] 连接 【PK服务失败】", zap.Error(err))
		return
	}
	global.PKClient = userpb.NewPKClient(conn)
	global.UserDivide = divide.NewDivide(global.PKClient, global.UserEnterSubscriber, global.UserCompleteSubscriber)
}

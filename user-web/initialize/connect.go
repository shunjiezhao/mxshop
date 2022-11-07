package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"web-api/user-web/global"
	userpb "web-api/user-web/proto"
)

// 连接grpc
func InitConnect() {
	// 连接 用户服务
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("[GetUserList] 连接 【用户服务失败】", zap.Error(err))
		return
	}
	global.UserServiceClient = userpb.NewUserServiceClient(conn)
}

package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
	"web-api/user-web/forms"
	"web-api/user-web/global"
	"web-api/user-web/global/response"
	userpb "web-api/user-web/proto"
)

func HandlerGrpcErrorToHttp(err error, c *gin.Context) {
	// 将 grpc code 转换为 http 状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用" + e.Message(),
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误" + e.Message(),
				})
			}
		}
	}
}

func GetUserList(ctx *gin.Context) {
	// 连接 用户服务
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("[GetUserList] 连接 【用户服务失败】", zap.Error(err))
		HandlerGrpcErrorToHttp(err, ctx)
		return
	}
	cli := userpb.NewUserServiceClient(conn)
	resp, err := cli.GetUserList(context.Background(), &userpb.PageInfo{
		Number: 1,
		Size:   0,
	})
	if err != nil {
		zap.L().Info("[GetUserList] 查询 【用户列表】失败")
		HandlerGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, val := range resp.Data {
		userResp := response.UserResponse{
			Id:       val.Id,
			NickName: val.NickName,
			Birthday: time.Unix(int64(val.Birthday), 0).Format("2022-11-26"),
			Gender:   val.Gender,
			Mobile:   val.Mobile,
		}

		//data["id"] = val.Id
		//data["name"] = val.NickName
		//data["birthday"] = val.Birthday
		//data["gender"] = val.Gender
		//data["mobile"] = val.Mobile
		result = append(result, userResp)
	}
	ctx.JSON(http.StatusOK, result)
}

func PassWordLogin(c *gin.Context) {
	//表单验证
	loginForm := forms.PassWordLoginForm{}
	if err, ok := forms.BindAndValid(c, &loginForm); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err,
		})
		return
	}

}

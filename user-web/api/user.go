package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
	"web-api/user-web/forms"
	"web-api/user-web/global"
	"web-api/user-web/global/response"
	userpb "web-api/user-web/proto"
	"web-api/user-web/utils/token"
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
	if checkGrpc(ctx) {
		return
	}
	claim, _ := ctx.Get("claim")
	customClaim := claim.(*token.CustomClaim)
	zap.L().Info("访问用户:", zap.String("name", customClaim.Nickname), zap.Int("id", int(customClaim.UserId)))
	resp, err := global.UserServiceClient.GetUserList(context.Background(), &userpb.PageInfo{
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

		result = append(result, userResp)
	}
	ctx.JSON(http.StatusOK, result)
}

func checkGrpc(ctx *gin.Context) bool {
	if global.UserServiceClient == nil {
		zap.L().Error("[GetUserList] 连接 【用户服务失败】")
		HandlerGrpcErrorToHttp(fmt.Errorf("连接 用户服务失败"), ctx)
		return true
	}
	return false
}

func PassWordLogin(c *gin.Context) {
	if checkGrpc(c) {
		return
	}
	//表单验证
	loginForm := forms.PassWordLoginForm{}
	if err, ok := forms.BindAndValid(c, &loginForm); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err,
		})
		return
	}
	user, err := global.UserServiceClient.GetUserByMobile(context.Background(), &userpb.GetUserByMobileRequest{
		Mobile: loginForm.Mobile,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			}
		}
		return
	}
	resp, err := global.UserServiceClient.CheckPassWord(context.Background(), &userpb.CheckPassWordRequest{
		PassWord: loginForm.PassWord,
		EncPwd:   user.PassWord,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
	}
	if resp != nil && resp.Success {
		token, err := global.JwtTokenGen.GenerateToken(user.NickName, user.Id, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "内部错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":    user.Id,
			"token": token,
			"msg":   "登陆成功",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "密码错误",
		})
	}

}

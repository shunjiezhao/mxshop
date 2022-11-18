package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
	"time"
	"web-api/user-web/forms"
	"web-api/user-web/global"
	"web-api/user-web/global/response"
	userpb "web-api/user-web/proto"
	"web-api/user-web/utils/token"
)

func RemoveTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}
func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": RemoveTopStruct(errs.Translate(global.Trans)),
	})
	return
}
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
	// 通过 etcd 获取 服务地址

	if global.ServerConfig.Debug == false {
		claim, _ := ctx.Get("claim")
		customClaim := claim.(*token.CustomClaim)
		zap.L().Info("访问用户:", zap.String("name", customClaim.Nickname), zap.Int("id", int(customClaim.UserId)))
	}
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
	if err, ok := forms.BindAndValid(c, &loginForm); !ok && !global.ServerConfig.Debug {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err,
		})
		return
	}
	if verify := store.Verify(loginForm.CaptchaId, loginForm.Captcha, true); !verify && !global.ServerConfig.Debug {
		c.JSON(http.StatusBadRequest, gin.H{"captcha": "验证码错误"})
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
		return
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
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "密码错误",
		})
		return
	}

}

// 用户注册
func Register(c *gin.Context) {
	if checkGrpc(c) {
		return
	}
	form := forms.RegisterForm{}
	if err, ok := forms.BindAndValid(c, &form); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err,
		})
		return
	}
	// 验证码校验
	value, err := global.Rdb.Get(context.Background(), form.Mobile).Result()
	if err != redis.Nil {
		zap.L().Info("用户注册的手机号码 没有发送验证码或者验证码过期")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码过期",
		})
		return
	}
	if form.Code != value {
		fmt.Println("want:%v; but:%v", value, form.Code)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	//注册
	user, err := global.UserServiceClient.CreateUser(context.Background(), &userpb.CreateUserRequest{
		Nickname: form.Mobile,
		Mobile:   form.Mobile,
		PassWord: form.PassWord,
	})
	if err != nil {
		zap.L().Info("[CreateUser]  【创建用户失败】", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{
			"msg": "无法创建用户",
		})
		return
	}
	token, err := global.JwtTokenGen.GenerateToken(user.NickName, user.Id, user.Role)
	c.JSON(http.StatusOK, gin.H{
		"id":    user.Id,
		"token": token,
	})
}

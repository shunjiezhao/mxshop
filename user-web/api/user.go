package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"time"
	"web-api/shared/queue"
	"web-api/shared/userid"
	"web-api/user-web/forms"
	"web-api/user-web/global"
	"web-api/user-web/global/response"
	userpb "web-api/user-web/proto"
	"web-api/user-web/utils/divide"
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

//TODO: 配置话
var testTime = time.Second * 5

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
		_, err = global.JWTTokenVerifier.Verify(token)
		if err != nil {
			fmt.Println("token err ", err)
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
	value, err := global.RedisClient.Get(context.Background(), form.Mobile).Result()
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

// 用户pk
func PK(c *gin.Context) {
	uid, err := userid.GetUid(c)
	if err != nil {
		HandlerGrpcErrorToHttp(err, c)
		return
	}
	tp, _ := strconv.Atoi(c.Param("type"))

	findType := userpb.FindType(tp)

	ch := make(chan *divide.Result)
	var clean func()
	recvMsg := make(chan []byte)

	// 在线匹配
	if findType == userpb.FindType_Random {
		// 接受匹配的通知
		if err := global.UserDivide.Register(queue.UserId(uid), ch, recvMsg); err != nil {
			// 已经注册
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		clean = func() {
			fmt.Println("开始clean")
			close(ch)
			close(recvMsg)
			global.UserDivide.UnRegister(queue.UserId(uid))
		}
	} else {
		defer close(ch)
		defer close(recvMsg)
	}

	resp, err := global.PKClient.Join(c, &userpb.JoinRequest{
		Id:       uid,
		FindType: findType,
		//TODO: 加入挑战的人的id
		OtherId: 0,
	})

	if err != nil {
		zap.L().Info("pk service 返回错误", zap.Error(err))
		HandlerGrpcErrorToHttp(err, c)
		return
	}
	fmt.Println(resp)
	// 在线匹配
	if findType == userpb.FindType_Random {
		result := <-ch // 获取到其他人的id
		fmt.Println("%d:获取到 result %v", result)
		OtherId := result.OtherID
		// 关闭管道 + 注销用户
		defer clean()
		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("建立websocket 连接失败", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		defer ws.Close()
		fmt.Println("建立完成")
		ctx, cancel := context.WithTimeout(context.Background(), testTime)
		defer cancel()

		// 建立 socket 连接 建立对局
		for {
			//读取ws中的数据
			select {
			case <-ctx.Done():
				// 对局结束
				// 发送分数 比较 并 发送结果给 gin
			default:
				// 写管道
				go func(recvMsg chan []byte) {
					fmt.Println(uid, "：开启监听websocket")
					for {
						select {
						case msg := <-recvMsg:
							//TODO: 可以写成抢答
							ws.WriteMessage(1, msg)
						case <-ctx.Done():
							return
						}
					}
					fmt.Println("end")
				}(recvMsg)
				_, message, err := ws.ReadMessage()
				if err != nil {
					return
				}
				// 写两次
				fmt.Println(uid, "->", OtherId, message, " ", recvMsg)

				global.UserDivide.SendMsg(OtherId, message)
				fmt.Println("send success")
			}
		}
	}
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func JoinParty(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uid, err := userid.GetUid(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	_, err = global.PKClient.TakePartIn(c, &userpb.TakePartInRequest{
		Id:  int32(id),
		Uid: uid,
	})
	if err != nil {
		HandlerGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

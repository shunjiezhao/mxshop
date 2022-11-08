package router

import (
	"github.com/gin-gonic/gin"
	"web-api/user-web/api"
)

func InitCaptcha(engine *gin.RouterGroup) {
	group := engine.Group("base")
	{
		group.GET("/captcha", api.GetDigCaptcha)
		group.POST("/sms", api.SendMessage) // 发送手机验证码
	}
}

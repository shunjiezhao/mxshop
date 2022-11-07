package router

import (
	"github.com/gin-gonic/gin"
	"web-api/user-web/api"
)

func InitCaptcha(engine *gin.RouterGroup) {
	group := engine.Group("captcha")
	{
		group.GET("/base64", api.GetDigCaptcha)
	}
}

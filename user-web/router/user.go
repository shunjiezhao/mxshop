package router

import (
	"github.com/gin-gonic/gin"
	"web-api/user-web/api"
	"web-api/user-web/middlewares"
)

func InitUserRouter(engine *gin.Engine) {
	uSrvPath := "/v1/user"
	// 登陆不需要jwttoken
	group := engine.Group(uSrvPath)
	{
		group.GET("list", middlewares.JwtToken(), api.GetUserList)
		group.POST("login", api.PassWordLogin)
	}
}

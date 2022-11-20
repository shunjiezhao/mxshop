package router

import (
	"github.com/gin-gonic/gin"
	"web-api/user-web/api"
	"web-api/user-web/global"
	"web-api/user-web/middlewares"
)

func InitUserRouter(engine *gin.RouterGroup) {
	// 登陆不需要jwttoken
	group := engine.Group("user")
	{
		if global.ServerConfig.Debug {
			group.GET("list", api.GetUserList)
		} else {
			group.GET("list", middlewares.JwtToken(), middlewares.IsAdmin(), api.GetUserList)
		}
		group.POST("login", api.PassWordLogin)
		group.POST("register", api.Register)
		group.GET("pk/:type", middlewares.JwtToken(), api.PK)
		// 参加活动 活动id
		group.GET("join/:id", middlewares.JwtToken(), api.JoinParty)

	}
}

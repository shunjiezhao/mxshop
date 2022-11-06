package router

import (
	"github.com/gin-gonic/gin"
	"web-api/user-web/api"
)

func InitUserRouter(engine *gin.Engine) {
	group := engine.Group("/v1/user")
	{
		group.GET("list", api.GetUserList)
		group.POST("login", api.PassWordLogin)
	}

}

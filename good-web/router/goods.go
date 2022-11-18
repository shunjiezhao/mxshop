package router

import (
	"github.com/gin-gonic/gin"
	"web-api/good-web/api/goods"
	"web-api/user-web/global"
)

func InitGoodsRouter(engine *gin.RouterGroup) {
	group := engine.Group("goods")
	{
		println(global.ServerConfig.Debug)
		group.GET("list", goods.GoodsList)
		group.POST("create", goods.CreateGoods) //改接口需要管理员权限
		//group.POST("create", middlewares.JwtToken(), middlewares.IsAdmin(), goods.CreateGoods) //改接口需要管理员权限
	}
}

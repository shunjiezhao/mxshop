package router

import (
	"github.com/gin-gonic/gin"
	"web-api/shop_cart-web/api/order"
	"web-api/shop_cart-web/global"
	"web-api/shop_cart-web/middlewares"
)

func InitOrderRouter(engine *gin.RouterGroup) {
	group := engine.Group("order").Use(middlewares.JwtToken())
	{
		println(global.ServerConfig.Debug)
		group.GET("list", middlewares.IsAdmin(), order.List) // 订单列表
		group.POST("", order.New)                            //新建订单
		group.GET("/:id", order.Detail)                      // 订单详情
		group.GET("pay/:ordersn", order.Pay)
	}
}

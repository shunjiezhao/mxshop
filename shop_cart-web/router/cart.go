package router

import (
	"github.com/gin-gonic/gin"
	"web-api/shop_cart-web/api/cart"
	"web-api/shop_cart-web/middlewares"
)

func InitCartRouter(engine *gin.RouterGroup) {
	group := engine.Group("cart").Use(middlewares.JwtToken())
	{
		group.GET("list", cart.List)      // 购物车列表
		group.POST("", cart.New)          //添加到购物车
		group.DELETE("/:id", cart.Delete) // 删除条目
		group.PATCH("/:id", cart.Update)  // 修改条目
	}
}

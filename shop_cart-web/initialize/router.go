package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"web-api/shop_cart-web/middlewares"
	"web-api/shop_cart-web/router"
	validator2 "web-api/shop_cart-web/validator"
)

func Routers() *gin.Engine {
	engine := gin.New()
	// 配置手机号的认证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", validator2.ValidateMobile)
	}
	// 配置跨域
	engine.Use(middlewares.Cors())
	apiGroup := engine.Group("/v1")
	{
		router.InitOrderRouter(apiGroup)
		router.InitCartRouter(apiGroup)
	}
	return engine
}

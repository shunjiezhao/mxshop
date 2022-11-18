package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"web-api/good-web/middlewares"
	"web-api/good-web/router"
	validator2 "web-api/good-web/validator"
)

func Routers() *gin.Engine {
	engine := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", validator2.ValidateMobile)
	}
	// 配置跨域
	engine.Use(middlewares.Cors())
	apiGroup := engine.Group("/v1")
	{
		router.InitGoodsRouter(apiGroup)

	}
	return engine
}

package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"web-api/user-web/router"
	validator2 "web-api/user-web/validator"
)

func Routers() *gin.Engine {
	engine := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", validator2.ValidateMobile)
	}
	router.InitUserRouter(engine)
	return engine
}

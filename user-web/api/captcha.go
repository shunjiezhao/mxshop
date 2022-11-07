package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

func GetDigCaptcha(c *gin.Context) {
	digit := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(digit, store)
	generate, s, err := captcha.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"base64": s,
		"id":     generate,
	})
}

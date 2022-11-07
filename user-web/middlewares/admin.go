package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web-api/user-web/utils/token"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, _ := c.Get("claim")
		customClaim := claim.(*token.CustomClaim)
		if customClaim.Role != 2 {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "没有权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web-api/user-web/global"
)

func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		if s, exist := c.GetQuery("token"); exist {
			token = s
		} else {
			token = c.GetHeader("token")
		}
		// 登陆不需要验证
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "没有携带token",
			})
		}
		claim, err := global.JWTTokenVerifier.Verify(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "token无效",
			})
			c.Abort()
		}
		c.Set("claim", claim)
		c.Set("user_id", claim.UserId)
		c.Next()
	}
}

package userid

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"web-api/shop_cart-web/api"
)

func GetUid(c *gin.Context) (int32, error) {
	uid, exists := c.Get("user_id")
	if !exists {
		api.HandleValidatorError(c, fmt.Errorf("user-id not exist"))
		return 0, fmt.Errorf("")
	}
	return uid.(int32), nil
}

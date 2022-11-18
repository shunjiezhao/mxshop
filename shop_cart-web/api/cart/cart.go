package cart

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	proto2 "web-api/good-web/proto"
	"web-api/shop_cart-web/api"
	"web-api/shop_cart-web/forms"
	"web-api/shop_cart-web/global"
	"web-api/shop_cart-web/proto"
)

func GetUid(c *gin.Context) (int32, error) {
	uid, exists := c.Get("user_id")
	if !exists {
		api.HandleValidatorError(c, fmt.Errorf("user-id not exist"))
		return 0, fmt.Errorf("")
	}
	return uid.(int32), nil
}
func List(c *gin.Context) {
	uid, err := GetUid(c)
	if err != nil {
		return
	}
	list, err := global.OrderClient.CartItemList(c, &proto.UserInfo{
		Id: uid,
	})
	if err != nil {
		api.HandlerGrpcErrorToHttp(err, c)
		return
	}
	if list.Total == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "空",
		})
		return
	}
	gids := make([]int32, list.Total)
	gid2Cidx := make(map[int32]int)
	for i, da := range list.Data {
		gid2Cidx[da.GoodsId] = i
		gids[i] = da.GoodsId
	}
	// 查询查品

	goods, err := global.GoodsClient.BatchGetGoods(c, &proto2.BatchGoodsIdInfo{Id: gids})
	if err != nil {
		api.HandlerGrpcErrorToHttp(err, c)
		return
	}
	type result struct {
		id         int32   `json:"id"`
		goodsId    int32   `json:"goods_id"`
		goodsName  string  `json:"goods_name"`
		goodsPrice float32 `json:"goods_price"`
		goodsImage string  `json:"goods_image"`
		nums       int32   `json:"nums"`
		checked    bool    `json:"checked"`
	}
	resp := gin.H{}
	resp["total"] = goods.Total
	goodsList := make([]interface{}, 0)
	for _, good := range goods.Data {
		cartD := list.Data[gid2Cidx[good.Id]] // 商品 id 对应购物车的id
		tmpMap := map[string]interface{}{}
		tmpMap["id"] = cartD.Id
		tmpMap["goods_id"] = good.Id
		tmpMap["good_name"] = good.Name
		tmpMap["good_image"] = good.GoodsFrontImage
		tmpMap["good_price"] = good.ShopPrice
		tmpMap["nums"] = cartD.Nums
		tmpMap["checked"] = cartD.Checked
		goodsList = append(goodsList, tmpMap)
	}

	resp["data"] = goodsList
	c.JSON(http.StatusOK, resp)
}
func Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := global.OrderClient.DeleteCartItem(c, &proto.CartItemRequest{
		Id: int32(id),
	})
	if err != nil {
		api.HandlerGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "OK"})

}
func New(c *gin.Context) {
	uid, err := GetUid(c)
	if err != nil {
		return
	}
	var form forms.CartItemForm
	if err := c.ShouldBindJSON(&form); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	cart, err := global.OrderClient.CreateCart(c, &proto.CartItemRequest{
		UserId:  uid,
		GoodsId: form.GoodsId,
		Nums:    form.Nums,
		Checked: false,
	})
	if err != nil {
		api.HandlerGrpcErrorToHttp(err, c)
		return
	}
	data := gin.H{}
	data["data"] = cart
	c.JSON(http.StatusOK, data)

}

func Update(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Param("id"))
	uid, err := GetUid(c)
	if err != nil {
		return
	}
	var form forms.UpdateCareItemForm
	if err := c.ShouldBindJSON(&form); err != nil {
		api.HandleValidatorError(c, err)
		return
	}
	if _, err := global.OrderClient.UpdateCartItem(c, &proto.CartItemRequest{
		UserId:  uid,
		GoodsId: int32(gid),
		Nums:    form.Nums,
		Checked: form.Checked,
	}); err != nil {
		api.HandlerGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "OK"})
}

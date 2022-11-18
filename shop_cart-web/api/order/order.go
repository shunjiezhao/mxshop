package order

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"web-api/good-web/utils/token"
	"web-api/shop_cart-web/api/cart"
	"web-api/shop_cart-web/forms"
	"web-api/shop_cart-web/global"
	"web-api/shop_cart-web/proto"
	"web-api/user-web/api"
)

func List(c *gin.Context) {
	uid, err := cart.GetUid(c)
	if err != nil {
		return
	}
	// 获取权限
	claims, _ := c.Get("claims")
	request := proto.OrderFilterRequest{}
	// 如果是管理员用户则返回所有的订单
	model := claims.(*token.CustomClaim)
	// 这是代表着普通用户
	if model.Role == 1 {
		request.UserId = uid
	}
	pages := c.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Page = int32(pagesInt)

	perNums := c.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)
	list, err := global.OrderClient.OrderList(c, &request)
	if err != nil {
		zap.S().Errorw("获取订单列表失败")
		api.HandlerGrpcErrorToHttp(err, c)
		return
	}

	reMap := gin.H{
		"total": list.Total,
	}
	orderList := make([]interface{}, 0)

	for _, item := range list.Data {
		tmpMap := map[string]interface{}{}
		tmpMap["id"] = item.Id
		tmpMap["status"] = item.Status
		tmpMap["pay_type"] = item.PayType
		tmpMap["user"] = item.UserId
		tmpMap["total"] = item.Total
		tmpMap["rcvInfo"] = item.RcvInfo
		tmpMap["order_sn"] = item.OrderSn
		tmpMap["id"] = item.Id
		orderList = append(orderList, tmpMap)
	}
	reMap["data"] = orderList
	c.JSON(http.StatusOK, reMap)
}
func New(c *gin.Context) {
	uid, err := cart.GetUid(c)
	if err != nil {
		return
	}
	var form forms.CreateOrderForm
	if err := c.ShouldBindJSON(&form); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	rsp, err := global.OrderClient.Create(c, &proto.OrderRequest{
		UserId: uid,
		RcvInfo: &proto.ReceiveInfo{
			Address: form.Address,
			RcvName: form.Name,
			Mobile:  form.Mobile,
			Post:    form.Post, // 备注
		},
	})
	if err != nil {
		zap.L().Info("无法创建订单", zap.Error(err))
		api.HandlerGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}
func Detail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "url参数错误",
		})
		return
	}

	uid, err := cart.GetUid(c)
	if err != nil {
		api.HandleValidatorError(c, err)
		return
	}
	rsp, err := global.OrderClient.OrderDetail(c, &proto.OrderRequest{
		Id:     int32(id),
		UserId: uid,
	})
	if err != nil {
		zap.S().Errorw("获取订单详情失败")
		api.HandlerGrpcErrorToHttp(err, c)
		return
	}
	reMap := gin.H{}
	reMap["id"] = rsp.OrderInfo.Id
	reMap["status"] = rsp.OrderInfo.Status
	reMap["user"] = rsp.OrderInfo.UserId
	reMap["total"] = rsp.OrderInfo.Total
	reMap["rcvinfo"] = rsp.OrderInfo.RcvInfo
	reMap["pay_type"] = rsp.OrderInfo.PayType
	reMap["order_sn"] = rsp.OrderInfo.OrderSn
	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Goods {
		tmpMap := gin.H{
			"id":    item.GoodsId,
			"name":  item.GoodName,
			"image": item.GoodsImage,
			"price": item.GoodsPrice,
			"nums":  item.Nums,
		}

		goodsList = append(goodsList, tmpMap)
	}
	reMap["goods"] = goodsList
	c.JSON(http.StatusOK, reMap)
}

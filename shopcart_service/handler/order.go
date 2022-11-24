package handler

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"math/rand"
	proto2 "server/goods_service/api/gen/v1/goods"
	proto3 "server/inventory_service/proto/gen/v1/inventory"
	"server/shopcart_service/global"
	"server/shopcart_service/model"
	proto "server/shopcart_service/proto/gen/v1/cart"
	"server/shopcart_service/utils/queue"
	"time"
)

type OrderService struct {
	db     *gorm.DB
	logger *zap.Logger
	GenId  func(int32) string // userid
	proto.UnsafeOrderServer
}

func (o *OrderService) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	var resp emptypb.Empty
	if o.db.Delete(&model.ShoppingCart{}, req.Id).RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "记录不存在")
	}
	return &resp, nil
}

type Config struct {
	DB     *gorm.DB
	Logger *zap.Logger
	GenId  func(int32) string
}

func New(config Config) *OrderService {
	val := &OrderService{
		db:     config.DB,
		logger: config.Logger,
	}
	if config.GenId != nil {
		val.GenId = config.GenId
	} else {
		val.GenId = defaultGenId
	}
	return val
}

//订单号生成 格式 年月日 用户id
func defaultGenId(id int32) string {
	now := time.Now()
	rand.Seed(now.UnixNano())
	return fmt.Sprintf("%d%d%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		id, rand.Intn(90)+10,
	)
}

func (o *OrderService) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var shopCarts []model.ShoppingCart
	result := o.db.Model(&model.ShoppingCart{User: req.Id}).Find(&shopCarts)
	if err := result.Error; err != nil {
		o.logger.Error("can not get shopcarts", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "")
	}

	rsp := &proto.CartItemListResponse{}
	rsp.Total = int32(result.RowsAffected)
	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}
	return rsp, nil
}

func (o *OrderService) CreateCart(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	var (
		cart model.ShoppingCart
		resp proto.ShopCartInfoResponse
	)
	result := o.db.Where(&model.ShoppingCart{User: req.UserId, Goods: req.GoodsId}).First(&cart)
	cart.Goods = req.GoodsId
	cart.User = req.UserId
	// 1. 购物车中没有这件商品，添加
	//2. 有数量 +1
	cart.Nums += req.Nums
	if result.RowsAffected == 0 {
		o.db.Create(&cart)
	} else {
		o.db.Model(&cart).Update("Nums", cart.Nums)
	}
	resp.Id = cart.ID
	resp.Nums = cart.Nums
	return &resp, nil
}

// 更新购物车 有几个操作
// 1. 更改数量 +-
// 2. 选中状态
// 3. 删除
func (o *OrderService) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	var (
		cart model.ShoppingCart
		resp emptypb.Empty
	)
	query := &model.ShoppingCart{User: req.UserId, Goods: req.GoodsId}
	result := o.db.Where(query).First(&cart)
	if result.RowsAffected == 0 {
		return &resp, status.Errorf(codes.NotFound, "没有该购物车信息")
	}
	cart.Nums += req.Nums      // 订单数量 这里需要注意 这里值得是我们 在原有的基础上增加多少次或减少多少次， 就像淘宝上面的 + - 我们只需要记录次数就可以了
	cart.Checked = req.Checked // 选中
	o.db.Save(&cart)
	return &resp, nil
}

func (o *OrderService) Create(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	// 1. 查询购物车中的商品
	// 2. 查询库存容量
	// 3. 查询商品金额
	// 4. 库存扣减
	// 1. 预扣减
	var order proto.OrderInfoResponse
	order.UserId = req.UserId
	var carts []model.ShoppingCart
	var goodsId []int32
	result := o.db.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Find(&carts)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "购物车中没有商品哦")
	}
	goods2Num := make(map[int32]int32)
	for _, cart := range carts {
		goodsId = append(goodsId, cart.Goods)
		goods2Num[cart.Goods] = cart.Nums // 方便后续进行库存比较
	}

	// 跨服务
	goodsList, err := global.GoodSrv.BatchGetGoods(ctx, &proto2.BatchGoodsIdInfo{Id: goodsId})
	if err != nil {
		o.logger.Info("无法从商品服务获取信息", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "")
	}
	var orderMount float32
	var orderGoods []*model.OrderGoods
	var invReq []*proto3.GoodsInvInfo
	for _, good := range goodsList.Data {
		cnt := goods2Num[good.Id]
		orderMount += good.ShopPrice * float32(cnt)
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       cnt,
		})
		invReq = append(invReq, &proto3.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     cnt,
		})
	}
	// 前面是所有的读
	orderModel := &model.OrderInfo{
		User:         req.UserId,
		OrderSn:      o.GenId(req.UserId),
		OrderMount:   orderMount,
		Address:      req.RcvInfo.Address,
		SignerName:   req.RcvInfo.RcvName,
		SingerMobile: req.RcvInfo.Mobile,
		Post:         req.RcvInfo.Post,
		Status:       model.PAYING,
	}
	sellInfo := proto3.SellInfo{
		GoodsInfo: invReq,
		OrderSn:   orderModel.OrderSn,
	}

	err = global.OrderPublisher.Publish(ctx, queue.OrderDelayQKey, &queue.OrderInfo{
		OrderSn:   orderModel.OrderSn,
		GoodsInfo: invReq,
		Test:      time.Now().String(),
	})
	if err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "创建订单失败")
	}

	_, err = global.InventorySrv.Sell(context.Background(), &sellInfo)
	if err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "扣减库存失败")
	}
	tx := o.db.Begin()
	// 订单号 时间桌
	tx.Save(orderModel)
	// 加入 orderGoods
	for _, orderGood := range orderGoods {
		orderGood.Order = orderModel.ID
	}
	// 批量插入
	fmt.Printf("%v\n", orderModel)
	if err := tx.CreateInBatches(&orderGoods, 1000).Error; err != nil {
		tx.Rollback()
		o.logger.Info("插入订单表商品表失败", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "")
	}
	// 删除购物车 勾选记录
	if err := tx.Unscoped().Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(&model.ShoppingCart{}).Error; err != nil {
		tx.Rollback()
		o.logger.Info("无法清除购物车状态", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "")
	}

	tx.Commit()
	resp := &proto.OrderInfoResponse{
		Id:      orderModel.ID, // 插入后订单的id
		UserId:  req.UserId,
		OrderSn: orderModel.OrderSn,
		Total:   order.Total,
		RcvInfo: order.RcvInfo,
	}
	return resp, nil
}

// 管理员 和 用户通用
func (o *OrderService) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderFilterResponse, error) {
	var orders []model.OrderInfo
	// 没有 用户值 就搜素全部
	var count int64
	o.db.Where(&model.OrderInfo{User: req.UserId}).Count(&count)
	var resp proto.OrderFilterResponse
	resp.Total = int32(count)
	o.db.Scopes(model.Paginate(int(req.Page), int(req.PagePerNums))).Find(&orders)
	for _, order := range orders {
		resp.Data = append(resp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Total:   order.OrderMount,
			RcvInfo: &proto.ReceiveInfo{
				Address: order.Address,
				RcvName: order.SignerName,
				Mobile:  order.SingerMobile,
			},
			Status: order.Status,
		})
	}
	return &resp, nil
}

// 传递的时候 也要传递 userid 用来判断 是否是该用户的订单，避免访问别人的
// 管理员查询所有的 利用gorm的零值的特性
func (o *OrderService) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	order := &model.OrderInfo{}
	order.ID = req.Id
	order.User = req.UserId
	result := o.db.Where(order).First(order)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有该订单信息")
	}
	orderInfo := proto.OrderInfoResponse{}
	orderInfo.Id = order.ID
	orderInfo.UserId = order.User
	orderInfo.OrderSn = order.OrderSn
	orderInfo.PayType = order.PayType
	orderInfo.Status = order.Status
	orderInfo.RcvInfo = &proto.ReceiveInfo{
		Address: order.Address,
		RcvName: order.SignerName,
		Mobile:  order.SingerMobile,
		Post:    order.Post,
	}
	var resp proto.OrderInfoDetailResponse
	resp.OrderInfo = &orderInfo
	// 查询订单的商品信息， 这因为这章表的出现，导致我们不用跨服务
	var orderGoods []model.OrderGoods
	o.db.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods)
	for _, orderGood := range orderGoods {
		resp.Goods = append(resp.Goods, &proto.OrderItemResponse{
			GoodsId:    orderGood.Goods,
			GoodName:   orderGood.GoodsName,
			GoodsPrice: orderGood.GoodsPrice,
			GoodsImage: orderGood.GoodsImage,
			Nums:       orderGood.Nums,
		})
	}
	return &resp, nil
}

// 当支付以后，徐更新订单状态
func (o *OrderService) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	// 条件更新
	if result := o.db.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); result.Error != nil || result.RowsAffected == 0 {
		o.logger.Info("订单状态更新失败")
		return nil, status.Errorf(codes.Internal, "更新失败")
	}
	return &emptypb.Empty{}, nil
}

func (o *OrderService) PayOrder(ctx context.Context, req *proto.PayOrderRequest) (*proto.PayOrderResponse, error) {
	err := global.OrderPublisher.Publish(ctx, queue.OrderFinishQKey, &queue.OrderInfo{
		OrderSn: req.OrderSn,
		Test:    time.Now().String(),
	})
	resp := &proto.PayOrderResponse{Msg: "OK"}
	if err != nil {
		resp.Msg = err.Error()
		return resp, err
	}
	return resp, nil
}

func (o *OrderService) Watch(ctx context.Context) {
	// 订单释放
	for {
		select {
		case <-ctx.Done():
			if _, ok := ctx.Value("close").(int); ok {
				close(global.OrderSubscriber.Finish)
				close(global.OrderSubscriber.Release)
			}
			return
		case info := <-global.OrderSubscriber.Finish:
			//完成支付
			o.logger.Info("支付订单成功", zap.String("ordersn", info.OrderSn))
			global.Rdb.Set(ctx, "order:finish:"+info.OrderSn, 1, time.Duration(queue.DelayOrderTimeMs)*time.Millisecond)
			// 删除未支付订单序列
			o.db.Where(&model.OrderInfo{OrderSn: info.OrderSn, Status: model.PAYING}).Select("Status",
				"IsDelete").Updates(model.OrderInfo{Status: model.TRADE_SUCCESS, BaseModel: model.BaseModel{IsDeleted: true}})
		case info := <-global.OrderSubscriber.Release:
			res, err := global.Rdb.Get(ctx, "order:finish:"+info.OrderSn).Result()
			if err != nil {
				o.logger.Info("得到key失败", zap.Error(err))
				break
			}
			if res == "1" {
				o.logger.Info("订单已经支付成功", zap.Error(err))
				break
			}
			// 释放订单
			o.logger.Info("释放订单", zap.String("ordersn", info.OrderSn))
			//TODO:利用Redis优化，查询订单ordersn是否未支付的消息队列里面
			result := o.db.Where(&model.OrderInfo{OrderSn: info.OrderSn, Status: "PAYING"}).Select("Status",
				"IsDelete").Updates(model.OrderInfo{Status: model.TRADE_CLOSED,
				BaseModel: model.BaseModel{IsDeleted: true}})
			fmt.Println(result)
		}
	}
}

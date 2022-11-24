package handler

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"server/inventory_service/global"
	"server/inventory_service/model"
	proto "server/inventory_service/proto/gen/v1/inventory"
	"server/inventory_service/utils/queue"
	"time"
)

//InventoryService
type InventoryService struct {
	logger     *zap.Logger
	db         *gorm.DB
	Publisher  *queue.Publisher
	Subscriber *queue.Subscriber
	proto.UnimplementedInventoryServer
}

type InventorySrvConfig struct {
	Logger     *zap.Logger
	DB         *gorm.DB
	Publisher  *queue.Publisher
	Subscriber *queue.Subscriber
}

func NewService(config *InventorySrvConfig) *InventoryService {
	return &InventoryService{
		db:         config.DB,
		logger:     config.Logger,
		Publisher:  config.Publisher,
		Subscriber: config.Subscriber,
	}
}

func (i *InventoryService) BatchInvDetail(ctx context.Context, req *proto.InvListRequest) (*proto.InvListResponse, error) {
	var inv []model.Inventory
	if i.db.Find(&inv, req.GoodsId).RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "没有库存信息")
	}
	var resp proto.InvListResponse
	resp.Total = int32(len(inv))
	for _, inventory := range inv {
		resp.Data = append(resp.Data, &proto.GoodsInvInfo{
			GoodsId: inventory.Goods,
			Num:     inventory.ID,
		})
	}
	return &resp, nil
}
func (i *InventoryService) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	// 没有新增
	var inv model.Inventory
	i.db.Where("goods = ?", req.GoodsId).Find(&inv)
	fmt.Println("%d", inv.Goods)
	inv.Goods = req.GoodsId // 没有增加 加入
	inv.Stocks = req.Num    // 设置库存

	if err := i.db.Save(&inv).Error; err != nil {
		i.logger.Error("can not save inventory", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	return &emptypb.Empty{}, nil
}

func (i *InventoryService) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if i.db.Where("goods = ? ", req.GoodsId).First(&inv).RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

func (i *InventoryService) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	//数据库基本的一个应用场景：数据库事务
	//并发情况之下 可能会出现超卖 1

	tx := i.db.Begin()
	//m.Lock() //获取锁 这把锁有问题吗？  假设有10w的并发， 这里并不是请求的同一件商品  这个锁就没有问题了吗？
	//这个时候应该先查询表，然后确定这个订单是否已经扣减过库存了，已经扣减过了就别扣减了
	//并发时候会有漏洞， 同一个时刻发送了重复了多次， 使用锁，分布式锁

	var details []model.GoodsDetail
	for _, goodInfo := range req.GoodsInfo {
		details = append(details, model.GoodsDetail{
			Goods: goodInfo.GoodsId,
			Num:   goodInfo.Num,
		})

		var inv model.Inventory

		mutex := global.RedisPool.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := tx.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}
	tx.Commit() // 需要自己手动提交操作
	// 传递自动解锁消息 避免订单系统宕机
	global.StockRebackPublisher.Publish(ctx, queue.StockReleaseQKey, queue.OrderInfo{
		OrderSn:   req.OrderSn,
		GoodsInfo: req.GoodsInfo,
		Test:      time.Now().String(),
	})
	return &emptypb.Empty{}, nil
}
func (i *InventoryService) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//TODO：加入缓存避免一直访问订单号 stock:reback:orderSn expire 30min

	// 1. 订单超时归还
	// 2. 订单创建失败，归还之前扣减的
	// 3. 手动归还
	// 扣减库存
	tx := i.db.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if i.db.Where("goods = ? ", goodInfo.GoodsId).First(&inv).RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Error(codes.NotFound, "没有库存信息")
		}
		// 扣减
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

//存在返回1
var judgeScript = `
if (redis.call('exists', KEYS[1]) == 1) then
    return 1
else
    redis.call('SET', KEYS[1], ARGV[1])
    redis.call('EXPIRE', KEYS[1], ARGV[2])
    return 0
end
`

func (i *InventoryService) WatchStockReback(ctx context.Context) {
	for {
		select {
		case msg := <-i.Subscriber.Release:
			// 如果返回 0 则表示没有释放过，那我们需要释放
			val, err := global.Rdb.Eval(ctx, judgeScript, []string{"order:AutoRelease:" + msg.OrderSn},
				1, 300).Result()
			if err != nil {
				i.logger.Info("脚本返回错误", zap.Error(err))
				break
			}
			fmt.Println(val, err)
			if val.(int64) == 1 {
				fmt.Println(val)
				i.logger.Info("库存已经被归还")
				break
			}
			i.logger.Info("得到释放库存的消息", zap.String("order_sn", msg.OrderSn))
			_, _ = i.Reback(ctx, &proto.SellInfo{
				GoodsInfo: msg.GoodsInfo,
				OrderSn:   msg.OrderSn,
			})
		}
	}

}

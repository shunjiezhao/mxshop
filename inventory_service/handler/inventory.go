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
	"server/inventory_service/proto"
)

type InventoryService struct {
	logger *zap.Logger
	db     *gorm.DB
	proto.UnimplementedInventoryServer
}

type InventorySrvConfig struct {
	Logger *zap.Logger
	DB     *gorm.DB
}

func NewService(config *InventorySrvConfig) proto.InventoryServer {
	return &InventoryService{
		db:     config.DB,
		logger: config.Logger,
	}
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
	return &emptypb.Empty{}, nil
}
func (i *InventoryService) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
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

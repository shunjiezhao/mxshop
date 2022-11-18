package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"server/goods_service/api/gen/v1/goods"
	"server/goods_service/model"
)

func goodsToResponse(goods *model.Goods) *proto.GoodsInfoResponse {
	return &proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

type GoodsServer struct {
	db *gorm.DB
	proto.UnimplementedGoodsServer
}

func (g *GoodsServer) GetGoodsListByIds(ctx context.Context, req *proto.GoodsListByIdsRequest) (*proto.GoodsListResponse, error) {
	if len(req.Id) == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "没有传递信息")
	}
	var (
		goods []model.Goods
		resp  proto.GoodsListResponse
	)

	result := g.db.Where("").Find(&goods, req.Id)
	if err := result.Error; err != nil {
		return nil, status.Errorf(codes.Internal, "")
	}
	for _, good := range goods {
		resp.Data = append(resp.Data, goodsToResponse(&good))
	}
	b, _ := json.Marshal(resp)
	println(string(b))
	return &resp, nil
}

func New(db *gorm.DB) *GoodsServer {
	return &GoodsServer{
		db: db,
	}
}

func (g *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	// 关键词搜素，查询新品，查询热门商品，通过价格区间筛选，通过商品分配筛选
	var resp proto.GoodsListResponse
	db := g.db.Model(&model.Goods{})

	if req.KeyWords != "" {
		db = db.Where("name like ?", "%"+req.KeyWords+"%")
	}
	if req.IsHot {
		db = db.Where("is_hot = true")
	}
	if req.IsNew {
		db = db.Where("is_new  = true")
	}
	if req.PriceMin > 0 {
		db = db.Where("shop_price >= ?", req.PriceMin)
	}
	if req.PriceMax > 0 {
		db = db.Where("shop_price <= ?", req.PriceMax)
	}

	if req.Brand > 0 {
		db = db.Where("brands_id = ?", req.Brand)
	}
	// 通过分类去查询
	// 1. 点击一级
	if req.TopCategory > 0 {
		var cate model.Category
		if g.db.Find(&cate, req.TopCategory).RowsAffected == 0 {
			return nil, status.Error(codes.NotFound, "没有标签")
		}
		var query string // 三级目录的查询串
		switch cate.Level {
		case 1:
			// 1. 找出二级目录
			// (select id from category where parent_category_id = %d) b
			// 2. 找出二级目录下的三级目录
			// select id from b where parent_category_id in b
			query = fmt.Sprintf("select id from categorys where parent_category_id in (select id from categorys where parent_category_id = %d)", cate.ID)
		case 2:
			// 找出三级目录来就行
			query = fmt.Sprintf("select id from categorys where parent_category_id = %d", cate.ID)
		case 3:
			query = fmt.Sprintf("select id from categorys where id = %d", cate.ID)
		}
		db = db.Where(fmt.Sprintf("category_id in (%s)", query))
	}
	var count int64
	db.Count(&count)
	resp.Total = int32(count)

	var goods []model.Goods
	db.Find(&goods)

	if req.Pages == 0 {
		req.Pages = 1
	}
	if req.PagePerNums == 0 {
		req.PagePerNums = 100
	}

	result := db.Preload("Category").Preload("Brands").Scopes(model.Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)
	if err := result.Error; err != nil {
		//zap.L().Info("can not get all goods", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	for _, good := range goods {
		resp.Data = append(resp.Data, goodsToResponse(&good))
	}
	b, _ := json.Marshal(resp)
	println(string(b))
	return &resp, nil
}

// 批量查询商品的信息
func (g *GoodsServer) BatchGetGoods(ctx context.Context, info *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	var (
		goods []model.Goods
		resp  proto.GoodsListResponse
	)

	result := g.db.Find(&goods, info.Id)
	if err := result.Error; err != nil {
		return nil, status.Error(codes.Internal, "无法获取商品")
	}
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.Internal, "空")
	}
	resp.Total = int32(len(goods))

	for _, good := range goods {
		resp.Data = append(resp.Data, goodsToResponse(&good))
	}
	return &resp, nil
}

func (g *GoodsServer) GetGoodsDetail(ctx context.Context, request *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods
	result := g.db.Preload("Category").Preload("Brands").Find(&goods, request.Id)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	resp := goodsToResponse(&goods)
	return resp, nil
}
func (g *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var resp proto.GoodsInfoResponse
	// 先判断 改category 和 brand 是否存在
	var cate model.Category
	if g.db.Find(&cate, req.CategoryId).RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "没有标签")
	}
	var brand model.Brands
	if g.db.Find(&brand, req.BrandId).RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "没有商家")
	}

	// 看商品是否存在
	//TODO：上传图片
	goods := model.Goods{
		Brands:          brand,
		BrandsID:        brand.ID,
		Category:        cate,
		CategoryID:      cate.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}

	//srv之间互相调用了
	result := g.db.Create(&goods)
	if err := result.Error; err != nil {
		zap.L().Info("can not create goods", zap.Error(err))
		return nil, status.Error(codes.Internal, "没有创建成功")
	}

	if result.RowsAffected == 0 {
		return nil, status.Error(codes.Internal, "没有创建成功")
	}
	resp.Id = int32(result.RowsAffected)
	return &resp, nil
}

func (g *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	result := g.db.Delete(&model.Goods{}, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "不存在商品")
	}
	return &emptypb.Empty{}, nil
}

func (g *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var goods model.Goods

	if result := g.db.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if result := g.db.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := g.db.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	//TODO：上传图片
	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	tx := g.db.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

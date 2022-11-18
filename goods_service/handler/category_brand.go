package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	proto "server/goods_service/api/gen/v1/goods"
	"server/goods_service/model"
)

func (g *GoodsServer) CategoryBrandList(ctx context.Context, req *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	var (
		resp    proto.CategoryBrandListResponse
		mre     []model.GoodsCategoryBrand
		records []*proto.CategoryBrandResponse
		total   int64
	)

	result := g.db.Model(&model.GoodsCategoryBrand{}).Count(&total)
	resp.Total = int32(total)
	if err := result.Error; err != nil {
		zap.L().Info("GoodsCategoryBrand : can not get count", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	g.db.Preload("Category").Preload("Brands").Scopes(model.Paginate(int(req.Pages), int(req.PagePerNums))).Find(&mre)
	for _, re := range mre {
		records = append(records, &proto.CategoryBrandResponse{
			Id:       re.ID,
			Brand:    brand2Info(&re.Brands),
			Category: cate2Info(&re.Category),
		})
	}
	resp.Data = records
	return &resp, nil
}

func (g *GoodsServer) GetCategoryBrandList(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	var (
		resp    proto.BrandListResponse
		records []*proto.BrandInfoResponse
		mcb     []model.GoodsCategoryBrand
	)
	g.db.Preload("Brands").Where("category_id = ?", req.Id).Find(&mcb)
	for _, brand := range mcb {
		records = append(records, brand2Info(&brand.Brands))
	}
	resp.Data = records
	resp.Total = int32(len(records))
	return &resp, nil
}

func (g *GoodsServer) CreateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	var (
		resp  proto.CategoryBrandResponse
		cate  model.Category
		brand model.Brands
	)

	if first := g.db.First(&cate, req.CategoryId); first.RowsAffected == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "没有目录")
	}
	resp.Category = cate2Info(&cate)
	if first := g.db.First(&brand, req.BrandId); first.RowsAffected == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "没有商家")
	}
	resp.Brand = brand2Info(&brand)

	insert := model.GoodsCategoryBrand{
		CategoryID: req.CategoryId,
		BrandsID:   req.BrandId,
	}
	if req.Id != 0 {
		insert.ID = req.Id
	}
	result := g.db.Create(insert)
	if err := result.Error; err != nil || result.RowsAffected == 0 {
		zap.L().Info("can not CreateCategoryBrand", zap.Error(err))
		return nil, status.Errorf(codes.FailedPrecondition, "")
	}
	resp.Id = int32(result.RowsAffected)
	return &resp, nil

}

func (g *GoodsServer) DeleteCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	result := g.db.Where("category_id = ? and brands_id = ?", req.CategoryId, req.BrandId).Delete(&model.GoodsCategoryBrand{})
	if result.Error != nil || result.RowsAffected == 0 {
		return nil, status.Error(codes.FailedPrecondition, "")
	}
	return &emptypb.Empty{}, nil
}

func (g *GoodsServer) UpdateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	// 如果修改后的 （c,b) 早已存在，那么我们直接将原来的删除即可
	// 如果不存在，更改原来的记录
	var (
		mcb   model.GoodsCategoryBrand
		cate  model.Category
		brand model.Brands
	)
	if first := g.db.First(&mcb, req.Id); first.RowsAffected == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "没有记录")
	}
	if first := g.db.First(&cate, req.CategoryId); first.RowsAffected == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "没有目录")
	}

	if first := g.db.First(&brand, req.BrandId); first.RowsAffected == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "没有商家")
	}
	if first := g.db.Where("category_id = ? and brands_id = ?", req.CategoryId, req.BrandId); first.RowsAffected != 0 {
		g.db.Where("category_id = ? and brands_id = ?", mcb.CategoryID, mcb.BrandsID).Delete(&model.GoodsCategoryBrand{})
		return &emptypb.Empty{}, nil
	}

	mcb.CategoryID = req.CategoryId
	mcb.BrandsID = req.BrandId
	g.db.Where("id = ?", req.Id).Updates(&mcb)
	return &emptypb.Empty{}, nil
}

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

func (g *GoodsServer) BannerList(ctx context.Context, empty *emptypb.Empty) (*proto.BannerListResponse, error) {
	var (
		resp          proto.BannerListResponse
		brandResponse []*proto.BannerResponse
		brands        []model.Banner
	)
	result := g.db.Find(&brands)
	resp.Total = int32(result.RowsAffected)
	for _, brand := range brands {
		brandResponse = append(brandResponse, banner2resp(&brand))
	}
	resp.Data = brandResponse
	return &resp, nil
}

func banner2resp(brand *model.Banner) *proto.BannerResponse {
	return &proto.BannerResponse{
		Id:    brand.ID,
		Image: brand.Image,
		Index: brand.Index,
		Url:   brand.Url,
	}
}

func req2Banner(req *proto.BannerRequest) *model.Banner {
	return &model.Banner{
		BaseModel: model.BaseModel{ID: req.Id},
		Image:     req.Image,
		Index:     req.Index,
		Url:       req.Url,
	}
}

func (g *GoodsServer) CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	banner := req2Banner(req)
	result := g.db.Create(banner)
	if err := result.Error; err != nil {
		zap.L().Info("can not create banner", zap.Int("id", int(req.Id)), zap.Error(err))
		return nil, status.Error(codes.FailedPrecondition, "")
	}
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.FailedPrecondition, "无法创建")
	}
	return banner2resp(banner), nil
}

func (g *GoodsServer) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	dest := &model.Banner{}
	if result := g.db.First(dest, req.Id); result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "没有改轮播图")
	}
	g.db.Delete(dest, req.Id)
	return &emptypb.Empty{}, nil

}

func (g *GoodsServer) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	brand := model.Banner{
		BaseModel: model.BaseModel{ID: req.Id},
		Image:     req.Image,
		Url:       req.Url,
	}
	// 注意零值不会更新
	resp := &emptypb.Empty{}
	result := g.db.Model(&brand).Where("id=?", req.Id).Updates(brand)
	if err := result.Error; err != nil {
		zap.L().Info("can not update brand", zap.Int("id", int(req.Id)), zap.Error(err))
		return resp, status.Error(codes.FailedPrecondition, "")
	}
	return resp, nil
}

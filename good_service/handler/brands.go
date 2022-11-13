package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	proto "server/good_service/api/gen/v1/goods"
	"server/good_service/model"
)

func (g *GoodsServer) BrandList(ctx context.Context, request *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	var (
		resp          proto.BrandListResponse
		brandResponse []*proto.BrandInfoResponse
		brands        []model.Brands
	)
	if g.db == nil {
		panic("DB is nil")
	}
	result := g.db.Find(&brands)
	resp.Total = int32(result.RowsAffected)
	for _, brand := range brands {
		brandResponse = append(brandResponse, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	resp.Data = brandResponse
	return &resp, nil
}

func (g *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	var (
		resp proto.BrandInfoResponse
	)

	result := g.db.Create(&model.Brands{
		BaseModel: model.BaseModel{ID: req.Id},
		Name:      req.Name,
		Logo:      req.Logo,
	})

	resp.Id = int32(result.RowsAffected)
	resp.Name = req.Name
	resp.Logo = req.Logo
	return &resp, nil
}

func (g *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	g.db.Delete(&model.Brands{}, req.Id)
	return &emptypb.Empty{}, nil
}

func (g *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brand := model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	resp := &emptypb.Empty{}
	result := g.db.Model(&brand).Where("id=?", req.Id).Updates(brand)
	if err := result.Error; err != nil {
		zap.L().Info("can not update brand", zap.Int("id", int(req.Id)), zap.Error(err))
		return resp, status.Error(codes.FailedPrecondition, "")
	}

	return resp, nil
}

func brand2Info(brands *model.Brands) *proto.BrandInfoResponse {
	return &proto.BrandInfoResponse{
		Id:   brands.ID,
		Name: brands.Name,
		Logo: brands.Logo,
	}
}

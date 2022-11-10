package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	proto "server/good_service/api/gen/v1/goods"
	"server/good_service/model"
	"server/user_service/global"
)

func (g *GoodsServer) BrandList(ctx context.Context, request *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	var (
		resp          proto.BrandListResponse
		brands        []model.Brands
		brandResponse []*proto.BrandInfoResponse
	)
	result := global.DB.Find(brands)
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

func (g *GoodsServer) CreateBrand(ctx context.Context, request *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) DeleteBrand(ctx context.Context, request *proto.BrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) UpdateBrand(ctx context.Context, request *proto.BrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

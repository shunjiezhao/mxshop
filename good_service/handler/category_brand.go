package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	proto "server/good_service/api/gen/v1/goods"
)

func (g *GoodsServer) CategoryBrandList(ctx context.Context, request *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) GetCategoryBrandList(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) CreateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) DeleteCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) UpdateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

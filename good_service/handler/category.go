package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	proto "server/good_service/api/gen/v1/goods"
)

func (g *GoodsServer) GetAllCategorysList(ctx context.Context, empty *emptypb.Empty) (*proto.CategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) GetSubCategory(ctx context.Context, request *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) CreateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) UpdateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

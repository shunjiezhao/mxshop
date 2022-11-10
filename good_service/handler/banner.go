package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	proto "server/good_service/api/gen/v1/goods"
)

func (g *GoodsServer) BannerList(ctx context.Context, empty *emptypb.Empty) (*proto.BannerListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) CreateBanner(ctx context.Context, request *proto.BannerRequest) (*proto.BannerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) DeleteBanner(ctx context.Context, request *proto.BannerRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) UpdateBanner(ctx context.Context, request *proto.BannerRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

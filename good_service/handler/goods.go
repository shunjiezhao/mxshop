package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"server/good_service/api/gen/v1/goods"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

func (g *GoodsServer) GoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) BatchGetGoods(ctx context.Context, info *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) CreateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) DeleteGoods(ctx context.Context, info *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) UpdateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) GetGoodsDetail(ctx context.Context, request *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

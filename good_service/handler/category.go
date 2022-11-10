package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	proto "server/good_service/api/gen/v1/goods"
	"server/good_service/global"
	"server/good_service/model"
)

func (g *GoodsServer) GetAllCategorysList(ctx context.Context, empty *emptypb.Empty) (*proto.CategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	var (
		resp  proto.SubCategoryListResponse
		cates []model.Category
	)
	// 找到子分类
	// 1. 找到parent_id = req.id
	find := global.DB.Find(&cates, "parent_category_id = ?", req.Id)
	resp.Total = int32(find.RowsAffected)
	if int(resp.Total) != len(cates) {
		log.Printf("GetSubCategory: len:%d; resp.Total:%d;", len(cates), resp.Total)
	}

	data := []*proto.CategoryInfoResponse{}
	for _, cate := range cates {
		data = append(data, cate2Info(&cate))
	}
	resp.SubCategorys = data
	return &resp, nil
}

func cate2Info(cate *model.Category) *proto.CategoryInfoResponse {
	return &proto.CategoryInfoResponse{
		Id:             cate.ID,
		Name:           cate.Name,
		ParentCategory: cate.ParentCategoryID,
		Level:          cate.Level,
		IsTab:          cate.IsTab,
	}
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

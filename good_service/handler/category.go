package handler

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	proto "server/good_service/api/gen/v1/goods"
	"server/good_service/model"
)

//1,,,,,1-a,1,false,
//2,,,,,2-a,2,false,1
//3,,,,,2-b,2,false,1
//4,,,,,3-a,3,false,2
//5,,,,,3-b,3,false,2

func (g *GoodsServer) GetAllCategorysList(ctx context.Context, empty *emptypb.Empty) (*proto.CategoryListResponse, error) {
	var (
		cates []model.Category
		resp  proto.CategoryListResponse
	)
	result := g.db.Model(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&cates)
	if err := result.Error; err != nil {
		return nil, status.Error(codes.NotFound, "")
	}
	b, err := json.Marshal(cates)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}
	resp.JsonData = string(b)

	return &resp, nil
}

// 子目录是 下子层目录
func (g *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	var (
		resp  proto.SubCategoryListResponse
		cates []model.Category
		cate  model.Category
	)
	if g.db.Find(&cate, "id = ? and level = ?", req.Id, req.Level).RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "不存在该目录")
	}

	resp.Info = cate2Info(&cate)
	// 找到子分类
	// 找到父节点是当前节点的
	find := g.db.Where(&model.Category{ParentCategoryID: req.Id}).Find(&cates)

	resp.Total = int32(len(cates) + int(find.RowsAffected))
	if int(resp.Total) != len(cates) {
		log.Printf("GetSubCategory: len:%d; resp.Total:%d;", len(cates), resp.Total)
	}

	data := make([]*proto.CategoryInfoResponse, resp.Total)
	for i, cate := range cates {
		data[i] = cate2Info(&cate)
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
func info2Cate(resp *proto.CategoryInfoRequest) *model.Category {
	return &model.Category{
		BaseModel:        model.BaseModel{ID: resp.Id},
		Name:             resp.Name,
		Level:            resp.Level,
		IsTab:            resp.IsTab,
		ParentCategoryID: resp.ParentCategory,
	}
}

func (g *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	var (
		cate model.Category
	)
	mp := make(map[string]interface{})
	if req.Id != 0 {
		mp["id"] = req.Id
	}
	if req.ParentCategory != 0 {
		if g.db.Find(&model.Category{}, req.ParentCategory).RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "父目录不存在")
		}
		mp["parent_category_id"] = req.ParentCategory
	}
	if req.Name != "" {
		mp["name"] = req.Name
	}

	mp["is_tab"] = req.IsTab
	if req.Level != 0 {
		mp["level"] = req.Level
	}
	cm := &model.Category{}
	if res := g.db.Find(cm, req.Id); res.RowsAffected != 0 {
		return nil, status.Error(codes.AlreadyExists, "标签已经存在")
	}

	result := g.db.Model(cm).Create(mp)
	if err := result.Error; err != nil {
		zap.L().Info("can not crate Category", zap.Error(err))
		return nil, status.Error(codes.FailedPrecondition, "")
	}
	if result.RowsAffected == 0 {
		zap.L().Info("can not crate Category")
		return nil, status.Error(codes.AlreadyExists, "")
	}
	zap.L().Info("create category success", zap.Any("category", cate))
	return cate2Info(&cate), nil
}

func (g *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	result := g.db.Delete(&model.Category{}, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "不存在category")
	}
	return &emptypb.Empty{}, nil
}

// 这里需要注意的是我们的 level > 0,
func (g *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	mp := make(map[string]interface{})
	if req.ParentCategory != 0 {
		mp["parent_category_id"] = req.ParentCategory
	}
	if req.Name != "" {
		mp["name"] = req.Name
	}
	mp["is_tab"] = req.IsTab
	if req.Level != 0 {
		mp["level"] = req.Level
	}
	// 零值也会更新哦
	result := g.db.Model(&model.Category{}).Where("id = ?", req.Id).Updates(mp)
	if err := result.Error; err != nil {
		zap.L().Info("can not Update Category ", zap.Error(err))
		return nil, status.Error(codes.FailedPrecondition, "")
	}
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "没有该cate哦")
	}
	return &emptypb.Empty{}, nil
}

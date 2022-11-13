package handler

import (
	"context"
	fmt "fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	proto "server/good_service/api/gen/v1/goods"
	"server/good_service/model"
	Potesting "server/shared/postgres/testing"
	"testing"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	db, err := Potesting.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	srv := GoodsServer{db: db}

	var mc model.Category
	srv.db.First(&mc, 1)
	insertData := []*model.Category{}
	var i int32 = 0
	for ; i < 3; i++ {
		insertData = append(insertData, &model.Category{
			BaseModel:        model.BaseModel{ID: i + 1},
			Name:             fmt.Sprintf("%d-%c", i+1, ('a' + i)),
			Level:            i + 1,
			IsTab:            i == 0,
			ParentCategoryID: i,
		})
		fmt.Printf("%v\n", insertData)
	}
	cases := []struct {
		name    string
		op      func() error
		wantErr bool
	}{
		{
			name: "getEmptyList",
			op: func() error {
				list, err := srv.GetAllCategorysList(ctx, &emptypb.Empty{})
				if err != nil {
					return err
				}
				if list.JsonData != "[]" {
					return fmt.Errorf("get error")
				}
				return nil
			},
		},
		{
			//TODO: 既然无法格式字符串定义结构来解析我们需要的字段
			name: "create_three",
			op: func() error {
				for _, data := range insertData {
					_, err := srv.CreateCategory(ctx, &proto.CategoryInfoRequest{
						Id:             data.ID,
						Name:           data.Name,
						ParentCategory: data.ParentCategoryID,
						Level:          data.Level,
						IsTab:          data.IsTab,
					})
					if err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			name: "get_three_insert_data",
			op: func() error {
				list, err := srv.GetAllCategorysList(ctx, &emptypb.Empty{})
				if err != nil {
					return err
				}
				wantString := `[{"id":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","name":"1-a","level":1,"is_tab":true,"parent_category_id":0,"sub_category":[{"id":2,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","name":"2-b","level":2,"is_tab":false,"parent_category_id":1,"sub_category":[{"id":3,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","name":"3-c","level":3,"is_tab":false,"parent_category_id":2,"sub_category":null}]}]},{"id":2,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","name":"2-b","level":2,"is_tab":false,"parent_category_id":1,"sub_category":[{"id":3,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","name":"3-c","level":3,"is_tab":false,"parent_category_id":2,"sub_category":[]}]},{"id":3,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","name":"3-c","level":3,"is_tab":false,"parent_category_id":2,"sub_category":[]}]`
				if list.JsonData != wantString {
					return fmt.Errorf("want:%v; but:%v\n", wantString, list.JsonData)
				}
				return nil
			},
		},
		{
			name: "exist_create",
			op: func() error {
				_, err := srv.CreateCategory(ctx, &proto.CategoryInfoRequest{
					Id:    4,
					Name:  "1-a",
					Level: 1,
				})
				if err != nil {
					return err
				}
				return nil
			},
			wantErr: true,
		},
		{
			name: "GetSubCategory",
			op: func() error {
				list, err := srv.GetSubCategory(ctx, &proto.CategoryListRequest{
					Id:    2,
					Level: 2,
				})
				if err != nil {
					return err
				}
				// level 1 + level 2
				if list.Total != 2 {
					return fmt.Errorf("total != data len, want:3; but:%d\n", list.Total)
				}
				if list.Info.Id != 2 {
					return fmt.Errorf("identity info is not equal")
				}
				if list.SubCategorys[0].Name != "3-c" {
					return fmt.Errorf("sub category name is not equal")
				}
				return nil
			},
		},
		{
			name: "update_zero_value",
			op: func() error {
				wantName := "test1"
				srv.UpdateCategory(ctx, &proto.CategoryInfoRequest{
					Id:    1,
					Name:  wantName,
					Level: 1,
					IsTab: false,
				})
				if err != nil {
					return err
				}
				var m model.Category
				db.First(&m, 1)
				if m.Name != wantName {
					return fmt.Errorf("want name:%q; but: %q;", wantName, m.Name)
				}
				if m.IsTab != false {
					return fmt.Errorf("don't update zero value")
				}
				return nil
			},
		},
		{
			name: "del_level1",
			op: func() error {
				_, err := srv.DeleteCategory(ctx, &proto.DeleteCategoryRequest{Id: 1})
				if err != nil {
					return err
				}
				if srv.db.Find(&model.Category{}, 1).RowsAffected != 0 {
					fmt.Errorf("delete fail")
				}
				return nil
			},
		},
		{
			name: "get_not_exist_sub_category",
			op: func() error {
				_, err := srv.GetSubCategory(ctx, &proto.CategoryListRequest{
					Id:    1,
					Level: 1,
				})
				return err
			},
			wantErr: true,
		},
	}

	var cates []model.Category
	for _, c := range cases {
		err := c.op()
		if c.wantErr {
			if err == nil {
				t.Errorf("%s: want err but not", c.name)
			} else {
				continue
			}
		}

		db.Find(&cates)
		fmt.Printf("%v\n", cates)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
		}
	}
}

func TestMain(m *testing.M) {
	Potesting.RunWithMongoInDocker(m)
}

func equal(req *proto.CategoryInfoRequest, resp *proto.CategoryInfoResponse) bool {
	if req.Id != resp.Id || req.Name != resp.Name || req.Level != resp.Level || req.IsTab != resp.IsTab || req.ParentCategory != resp.ParentCategory {
		return false
	}
	return true
}

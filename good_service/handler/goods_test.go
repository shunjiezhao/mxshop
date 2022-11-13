package handler

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"gorm.io/gorm"
	proto "server/good_service/api/gen/v1/goods"
	"server/good_service/model"
	Potesting "server/shared/postgres/testing"
	"testing"
)

var (
	brandsName   = []string{"b1", "b2", "b3", "b4"}
	cateName     = []string{"1-a", "2-a", "3-a", "1-b", "2-b", "3-b"}
	cateParentID = []int32{0, 1, 2, 0, 4, 5}
	goodsName    = []string{"g1a", "g2a", "g3", "g4"}
	goodsCateID  = []int32{3, 3, 6, 6}
)

func insertHelp() *gorm.DB {
	var i, j int32
	db, err := Potesting.NewClient(context.Background())
	if err != nil {
		panic(err)
	}
	// 创建 目录 1-a 2-a 3-a
	//    1-b 2-b 3-b

	for i = 0; i < 2; i++ {
		for j = 1; j <= 3; j++ {
			idx := i*3 + j
			mp := map[string]interface{}{
				"id":    idx,
				"name":  cateName[idx-1],
				"level": j,
			}
			if cateParentID[idx-1] != 0 {
				mp["parent_category_id"] = cateParentID[idx-1]
			}
			db.Model(&model.Category{}).Create(mp)
		}
	}
	// 创建商标
	var brands []model.Brands

	for i = 1; i <= 4; i++ {
		brands = append(brands, model.Brands{
			BaseModel: model.BaseModel{ID: i},
			Name:      brandsName[i-1],
		})
	}
	db.CreateInBatches(brands, 100)

	// 创建商标目录表
	// (3-a,b1) (3-a,b2) (3-b,b3) (3-b, b4)
	// 3,1 3,2 6,3 6,4
	var c2b []model.GoodsCategoryBrand
	for i = 1; i <= 2; i++ {
		c2b = append(c2b, model.GoodsCategoryBrand{
			CategoryID: 3,
			BrandsID:   i,
		})
		c2b = append(c2b, model.GoodsCategoryBrand{
			CategoryID: 6,
			BrandsID:   i + 2,
		})
	}
	db.CreateInBatches(c2b, 100)

	// 创建商品
	var goods []model.Goods
	//goodsBrandID := []int32{1,2,3,4}
	for i = 1; i <= 4; i++ {
		goods = append(goods, model.Goods{
			BaseModel:  model.BaseModel{ID: i},
			CategoryID: goodsCateID[i-1],
			BrandsID:   i,
			Name:       goodsName[i-1],
			ShopPrice:  float32(i * 10),
		})
	}
	db.CreateInBatches(goods, 100)
	return db
}
func TestGetGoods(t *testing.T) {
	db := insertHelp()
	srv := GoodsServer{db: db}
	ctx := context.Background()
	cases := []struct {
		name    string
		op      func() error
		wantErr bool
	}{
		// 价格区间
		{
			name: "价格区间 [0,10]",
			op: func() error {
				list, err := srv.GoodsList(ctx, &proto.GoodsFilterRequest{
					PriceMin: 0,
					PriceMax: 10,
				})
				if err != nil {
					return err
				}
				// get g1a
				if list.Total != 1 || list.Data[0].Name != "g1a" || list.Data[0].ShopPrice != 10 {
					return fmt.Errorf("get [0,10] goods fail")
				}
				return nil
			},
		},
		{
			name: "价格区间 [20,30]",
			op: func() error {
				list, err := srv.GoodsList(ctx, &proto.GoodsFilterRequest{
					PriceMin: 20,
					PriceMax: 30,
				})
				if err != nil {
					return err
				}
				// get g1a
				if list.Total != 2 || list.Data[0].Name != "g2a" || list.Data[0].ShopPrice != 20 ||
					list.Data[1].Name != "g3" || list.Data[1].ShopPrice != 30 {
					return fmt.Errorf("get [20,30] goods fail")
				}
				return nil
			},
		},
		{
			name: "价格区间 [40, ...]",
			op: func() error {
				list, err := srv.GoodsList(ctx, &proto.GoodsFilterRequest{
					PriceMin: 40,
				})
				if err != nil {
					return err
				}
				// get g1a
				if list.Total != 1 || list.Data[0].Name != "g4" || list.Data[0].ShopPrice != 40 {
					return fmt.Errorf("get [40, ...] goods fail")
				}
				return nil
			},
		},
		// 关键字
		{
			name: "关键字",
			op: func() error {
				list, err := srv.GoodsList(ctx, &proto.GoodsFilterRequest{KeyWords: "a"})
				if err != nil {
					return err
				}
				if list.Total != 2 || list.Data[0].Name != "g1a" || list.Data[1].Name != "g2a" {
					return fmt.Errorf("查询关键词失败")
				}
				return nil
			},
		},
		{
			name: "目录",
			op: func() error {
				var result []*proto.GoodsListResponse
				// level1 目录
				// level2 目录
				// level3 目录
				for i := 0; i < 3; i++ {
					list, err := srv.GoodsList(ctx, &proto.GoodsFilterRequest{
						TopCategory: int32(i + 1),
					})
					if err != nil {
						return err
					}
					result = append(result, list)
				}

				for i := 0; i < 3; i++ {
					if diff := cmp.Diff(result[i%3], result[(i+1)%3], protocmp.Transform()); diff != "" {
						return fmt.Errorf(diff)
					}
				}
				return nil
			},
		},
		{
			name: "聚合查询key和价格",
			op: func() error {
				list, err := srv.GoodsList(ctx, &proto.GoodsFilterRequest{
					KeyWords: "a",
					PriceMin: 11,
				})
				if err != nil {
					return err
				}
				if list.Total != 1 || list.Data[0].Name != "g2a" {
					return fmt.Errorf("list.Total: %d list.Data[0].Name %x", list.Total, list.Data[0].Name)
				}
				return nil
			},
		},
		{
			name: "聚合查询key和价格 但没有",
			op: func() error {
				list, err := srv.GoodsList(ctx, &proto.GoodsFilterRequest{
					KeyWords: "a",
					PriceMin: 40,
				})
				if err != nil {
					return err
				}
				if list.Total != 0 {
					return fmt.Errorf("want nil; but get: %v", list.Data)
				}
				return nil
			},
		},
		{
			name: "聚合查询目录和价格",
			op: func() error {
				list, err := srv.GoodsList(ctx, &proto.GoodsFilterRequest{
					TopCategory: 6,
					PriceMin:    20,
				})
				if err != nil {
					return err
				}

				if list.Total != 2 || list.Data[0].Name != "g3" || list.Data[1].Name != "g4" {
					return fmt.Errorf(" list.Data[0].Name want:g3;but:%s\n  list.Data[1].Name want:g4;but:%s\n",
						list.Data[0].Name, list.Data[1].Name)
				}
				return err
			},
		},
	}
	for _, cc := range cases {
		err := cc.op()
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: want err but not", cc.name)
			} else {
				continue
			}
		}
		if err != nil {
			t.Errorf("%s: %v", cc.name, err)
		}
	}

}

func TestGetGoodsDetail(t *testing.T) {
	db := insertHelp()
	srv := GoodsServer{db: db}
	ctx := context.Background()
	cases := []struct {
		name    string
		op      func() error
		wantErr bool
	}{
		{
			name: "get_exist",
			op: func() error {
				detail, err := srv.GetGoodsDetail(ctx, &proto.GoodInfoRequest{Id: 1})
				if err != nil {
					return err
				}
				if detail.Id != 1 || detail.Name != "g1a" || detail.ShopPrice != 10 {
					return fmt.Errorf("want: id:1 name:g1a price:10\n but:id:%d name:%s price%v\n", detail.Id, detail.Name, detail.ShopPrice)
				}
				return nil
			},
		},
		{
			name: "get_not_exist",
			op: func() error {
				_, err := srv.GetGoodsDetail(ctx, &proto.GoodInfoRequest{Id: 5})
				return err
			},
			wantErr: true,
		},
	}
	for _, cc := range cases {
		err := cc.op()
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: want err but not", cc.name)
			} else {
				continue
			}
		}
		if err != nil {
			t.Errorf("%s: %v", cc.name, err)
		}
	}
}
func TestBatchGoods(t *testing.T) {
	db := insertHelp()
	srv := GoodsServer{db: db}
	ctx := context.Background()
	cases := []struct {
		name    string
		op      func() error
		wantErr bool
	}{
		{
			name: "get",
			op: func() error {
				id := []int32{2, 3}
				res, err := srv.BatchGetGoods(ctx, &proto.BatchGoodsIdInfo{Id: id})
				if err != nil {
					return err
				}
				if res.Total != 2 {
					return fmt.Errorf("want len:2 but len:%d\n", res.Total)
				}
				for i := 0; i < 2; i++ {
					if goodsName[i+1] != res.Data[i].Name {
						return fmt.Errorf("want name:%s but len:%s\n", goodsName[i+1], res.Data[i].Name)
					}
					if goodsCateID[i+1] != res.Data[i].CategoryId {
						return fmt.Errorf("want categoryId:%d but len:%d\n", goodsCateID[i+1], res.Data[i].CategoryId)
					}
				}
				return nil
			},
		},
		{
			name: "get_not_exist",
			op: func() error {
				_, err := srv.BatchGetGoods(ctx, &proto.BatchGoodsIdInfo{Id: []int32{6}})
				return err
			},
			wantErr: true,
		},
	}
	for _, cc := range cases {
		err := cc.op()
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: want err but not", cc.name)
			} else {
				continue
			}
		}
		if err != nil {
			t.Errorf("%s: %v", cc.name, err)
		}
	}
}

package model

import (
	"github.com/lib/pq"
)

const (
	UserMobileFieldName = "mobile"
	IDFieldName         = "id"
)

type Category struct {
	BaseModel
	Name  string `gorm:"type:varchar(20);not null"`
	Level int32  `gorm:"type:int;not null;default:1"`
	IsTab bool   `gorm:"type:bool;not null;default:false"`

	ParentCategoryID int32
	ParentCategory   *Category
}

func (b *Category) TableName() string {
	return "categorys"
}

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}

func (b *Brands) TableName() string {
	return "brands"
}

type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique;"`
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique;"`
	Brands   Brands
}

func (b *GoodsCategoryBrand) TableName() string {
	return "goods_category_brand"
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);default:1;not null"`
}

func (b *Banner) TableName() string {
	return "banners"
}

// 两个晚间
type Goods struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category

	BrandsID int32 `gorm:"type:int;not null"`
	Brands   Brands
	OnSale   bool `gorm:"default:false;not null"` // 是否上架
	ShipFree bool `gorm:"default:false;not null"` // 包邮？
	IsNew    bool `gorm:"default:false;not null"` // 是新品嘛
	IsHot    bool `gorm:"default:false;not null"` // 是否热销 付费

	Name   string `gorm:"type:varchar(200);not null"` // 商品名称
	GoodSn string `gorm:"type:varchar(200);not null"` // 商品id 对于商家而言

	ClickNum int32 `gorm:"type:int;default:0;not null"` // 点击数量
	SoldNum  int32 `gorm:"type:int;default:0;not null"` //卖出数量
	FavNum   int32 `gorm:"type:int;default:0;not null"` // 收藏数量

	MarketPrice    float32        `gorm:"not null"`                   //市场价 必填
	ShopPrice      float32        `gorm:"not null"`                   //卖的价 必填
	GoodBrief      string         `gorm:"type:varchar(100);not null"` // 商品描述  必填
	Images         pq.StringArray `gorm:"type:text[]"`                //左侧轮播图
	DescImages     pq.StringArray `gorm:"type:text[]"`                // 详细描述的图片
	GoodFrontImage string         //封面图
}

func (b *Goods) TableName() string {
	return "goods"
}

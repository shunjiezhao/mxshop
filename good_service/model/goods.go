package model

import (
	"github.com/lib/pq"
)

const (
	UserMobileFieldName = "mobile"
	IDFieldName         = "id"
)

//1,,,,,1-a,1,false,
//2,,,,,2-a,2,false,1
//3,,,,,2-b,2,false,1
//4,,,,,3-a,3,false,2
//5,,,,,3-b,3,false,2
//6,,,,,1-b,1,false,
type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);not null;unique" json:"name"`
	Level            int32       `gorm:"type:int;not null;default:1" json:"level"` // 从一开始
	IsTab            bool        `gorm:"type:bool;not null;default:false" json:"is_tab"`
	ParentCategoryID int32       `json:"parent_category_id"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
}

func (b *Category) TableName() string {
	return "categorys"
}

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null" json:"name"`
	Logo string `gorm:"type:varchar(200);default:'';not null" json:"logo"`
}

func (b *Brands) TableName() string {
	return "brands"
}

type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32    `gorm:"type:int;index:idx_category_brand;" json:"-"`
	Category   Category `json:"category"`

	BrandsID int32  `gorm:"type:int;index:idx_category_brand;" `
	Brands   Brands `json:"brands"`
}

func (b *GoodsCategoryBrand) TableName() string {
	return "goods_category_brand"
}

type Banner struct {
	BaseModel
	Index int32  `gorm:"type:int;not null;default:0" json:"index"`
	Image string `gorm:"type:varchar(200);not null" json:"image"`
	Url   string `gorm:"type:varchar(200);default:1;not null" json:"url"`
}

func (b *Banner) TableName() string {
	return "banners"
}

// 两个晚间
type Goods struct {
	BaseModel

	CategoryID int32    `gorm:"type:int;not null" json:"-"`
	Category   Category `json:"category"`

	BrandsID int32  `gorm:"type:int;not null" json:"brands_id"`
	Brands   Brands `json:"brands"`
	OnSale   bool   `gorm:"default:false;not null" json:"on_sale"`   // 是否上架
	ShipFree bool   `gorm:"default:false;not null" json:"ship_free"` // 包邮？
	IsNew    bool   `gorm:"default:false;not null" json:"is_new" `   // 是新品嘛
	IsHot    bool   `gorm:"default:false;not null" json:"is_hot" `   // 是否热销 付费

	Name    string `gorm:"type:varchar(200);not null" json:"name"`               // 商品名称
	GoodsSn string `gorm:"type:varchar(200);not null;default:''" json:"good_sn"` // 商品id 对于商家而言

	ClickNum int32 `gorm:"type:int;default:0;not null" json:"click_num"` // 点击数量
	SoldNum  int32 `gorm:"type:int;default:0;not null" json:"sold_num"`  //卖出数量
	FavNum   int32 `gorm:"type:int;default:0;not null"`                  // 收藏数量

	MarketPrice     float32        `gorm:"not null"`                   //市场价 必填
	ShopPrice       float32        `gorm:"not null"`                   //卖的价 必填
	GoodsBrief      string         `gorm:"type:varchar(100);not null"` // 商品描述  必填
	Images          pq.StringArray `gorm:"type:text[]"`                //左侧轮播图
	DescImages      pq.StringArray `gorm:"type:text[]"`                // 详细描述的图片
	GoodsFrontImage string         //封面图
}

func (b *Goods) TableName() string {
	return "goods"
}

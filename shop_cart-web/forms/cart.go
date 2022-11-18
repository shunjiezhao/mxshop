package forms

type CartItemForm struct {
	GoodsId int32 `json:"goods_id" form:"goods_id" binding:"required,min=1"`
	Nums    int32 `json:"nums" form:"nums" binding:"required,min=1"`
}

type UpdateCareItemForm struct {
	Checked bool  `json:"checked" form:"checked"`
	Nums    int32 `json:"nums" form:"nums" binding:"min=0"`
}

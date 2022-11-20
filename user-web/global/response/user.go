package response

type UserResponse struct {
	Id       int32  `json:"id"`
	NickName string `json:"name"`
	Birthday string `json:"birthday"`
	Gender   uint32 `json:"gender"`
	Mobile   string `json:"mobile"`
}

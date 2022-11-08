package forms

type SendSmsForm struct {
	Mobile string `json:"mobile" form:"mobile" binding:"required,mobile"`
	// 1- register 2 -login
	Type int `json:"type" form:"type"binding:"required,oneof=1 2"` // 1. 注册发送短信验证码和动态发送验证嘛
}

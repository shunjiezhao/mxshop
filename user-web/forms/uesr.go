package forms

type PassWordLoginForm struct {
	Mobile    string `json:"mobile" form:"mobile" binding:"required,mobile"`
	PassWord  string `json:"password" form:"password" binding:"required,min=5,max=12"`
	Captcha   string `json:"captcha" form:"captcha" binding:"required"`
	CaptchaId string `json:"captcha_id" form:"captcha_id" binding:"required"`
}
type RegisterForm struct {
	Mobile   string `json:"mobile" form:"mobile" binding:"required,mobile"`
	PassWord string `json:"password" form:"password" binding:"required,min=5,max=12"`
	Code     string `json:"code" form:"code" binding:"required,max=4"`
}

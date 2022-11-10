package global

import "errors"

var (
	AuthenticatedErr     = errors.New("验证失败")
	UserNotExist         = errors.New("用户不存在")
	UserAlreadyExist     = errors.New("用户已经存在")
	NickNameAlreadyExist = errors.New("昵称已经存在")
	PassWordNotVerify    = errors.New("密码错误")
)

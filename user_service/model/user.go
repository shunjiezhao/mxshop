package model

import (
	"time"
)

const (
	UserMobileFieldName = "mobile"
	IDFieldName         = "id"
)

type User struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null"` // 建立所以
	PassWord string     `gorm:"type:varchar(100);not null"`
	NickName string     `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:DATE"`                                          // 防止保存出错 指针类型空值 为 Nil
	Gender   int32      `gorm:"type:int;default:1;check:gender<2;comment:0-女,1-男;"` // 0-女 1-男
	Role     int32      `gorm:"type:int;default:1;check:role<3;comment:1-common_user,2-admin;"`
}

func (u User) TableName() string {
	return "user"
}

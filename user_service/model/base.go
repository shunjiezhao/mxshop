package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32     `gorm:"primarykey;type:int"`
	CreatedAt time.Time `gorm:"comment:创建时间;autoCreateTime"`
	UpdatedAt time.Time `gorm:"comment:修改时间;autoUpdateTime"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

func (b BaseModel) TableName() string {
	return "base_model"
}

package dao

import (
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

type Dao struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Dao {
	return &Dao{DB: db}
}

type ErrResult struct {
	Code codes.Code
	Err  error
}

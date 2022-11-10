package dao

import (
	context "context"
	"errors"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
	userpb "server/user_service/api/gen/v1/user"
	"server/user_service/global"
	"server/user_service/model"
	"time"
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

func dbUserToDao(user *model.User) *userpb.UserInfo {
	info := &userpb.UserInfo{
		Id:       user.ID,
		Mobile:   user.Mobile,
		PassWord: user.PassWord, // 需要传输嘛
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     user.Role,
	}
	if user.Birthday != nil {
		info.Birthday = uint64(user.Birthday.Unix())
	}
	return info
}

func daoUserToDB(user *userpb.UserInfo) *model.User {
	info := &model.User{
		BaseModel: model.BaseModel{ID: user.Id},
		Mobile:    user.Mobile,
		PassWord:  user.PassWord, // 需要传输嘛
		NickName:  user.NickName,
		Gender:    user.Gender,
		Role:      user.Role,
	}
	if user.Birthday != 0 {
		ti := time.Unix(int64(user.Birthday), 0)
		info.Birthday = &ti
	}
	return info

}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// 所有用户
func (d *Dao) GetUserList(ctx context.Context, page, pageSize uint32) ([]*userpb.UserInfo, *ErrResult) {
	users := []model.User{}

	result := d.DB.Find(&users)
	if result.Error != nil {
		return nil, &ErrResult{
			Code: codes.NotFound,
			Err:  global.UserNotExist,
		}
	}
	var records []*userpb.UserInfo
	for _, user := range users {
		records = append(records, dbUserToDao(&user))
	}
	global.DB.Scopes(Paginate(int(page), int(pageSize))).Find(&users)

	return records, nil
}

func (d *Dao) GetUserByUnique(ctx context.Context, m map[string]interface{}) (*userpb.UserInfo, *ErrResult) {
	var user model.User
	result := d.DB.Where(m).First(&user)

	if result.RowsAffected == 0 || errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, &ErrResult{
			Code: codes.NotFound,
			Err:  global.UserNotExist,
		}
	}
	if err := result.Error; err != nil {
		return nil, &ErrResult{
			Code: codes.Internal,
			Err:  err,
		}
	}
	return dbUserToDao(&user), nil
}

func (d *Dao) CreateUser(ctx context.Context, user *userpb.UserInfo) *ErrResult {
	result := d.DB.Create(daoUserToDB(user))
	err := result.Error
	if errors.Is(err, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return &ErrResult{
			Code: codes.AlreadyExists,
			Err:  global.UserAlreadyExist,
		}
	}
	if err != nil {
		return &ErrResult{
			Code: codes.Internal,
			Err:  err,
		}
	}
	return nil
}

func (d *Dao) UpdateUser(ctx context.Context, info *userpb.UserInfo) *ErrResult {
	// 不存在
	var user model.User
	result := d.DB.Find(&user, info.Id)
	checkErr := func(db *gorm.DB) *ErrResult {
		if result.RowsAffected == 0 {
			return &ErrResult{
				Code: codes.NotFound,
				Err:  global.UserNotExist,
			}
		}
		return nil
	}
	if err := checkErr(result); err != nil {
		return err
	}
	user.ID = info.Id
	ti := time.Unix(int64(info.Birthday), 0)
	user.Birthday = &ti
	user.Gender = info.Gender
	result = d.DB.Save(&user)
	if err := checkErr(result); err != nil {
		return err
	}
	return nil
}

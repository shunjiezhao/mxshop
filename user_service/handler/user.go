package handler

import (
	context "context"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"server/user_service/api/gen/v1/user"
	"server/user_service/dao"
	"server/user_service/global"
	"server/user_service/model"
	"strings"
)

type UserService struct {
	Dao *dao.Dao
	userpb.UnimplementedUserServiceServer
}

// 获取用户列表
func (u *UserService) GetUserList(ctx context.Context, pageInfo *userpb.PageInfo) (*userpb.GetUserListResponse, error) {
	userList, err := u.Dao.GetUserList(ctx, pageInfo.Number, pageInfo.Size)

	if err != nil {
		return nil, status.Error(err.Code, err.Err.Error())
	}
	var resp userpb.GetUserListResponse
	resp.Total = int32(len(userList))
	// 分页 利用 limit
	for _, user := range userList {
		resp.Data = append(resp.Data, (user))
	}
	return &resp, nil
}

func (u *UserService) GetUserByMobile(ctx context.Context, req *userpb.GetUserByMobileRequest) (*userpb.UserInfo, error) {
	user, err := u.Dao.GetUserByUnique(ctx, map[string]interface{}{
		model.UserMobileFieldName: req.Mobile,
	})
	if err != nil {
		return nil, status.Error(err.Code, err.Err.Error())
	}
	return (user), nil
}

func (u *UserService) GetUserById(ctx context.Context, req *userpb.GetUserByIdRequest) (*userpb.UserInfo, error) {
	user, err := u.Dao.GetUserByUnique(ctx, map[string]interface{}{
		model.IDFieldName: req.Id,
	})
	if err != nil {
		return nil, status.Error(err.Code, err.Err.Error())
	}
	return (user), nil
}

func (u *UserService) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserInfo, error) {
	if req.PassWord == "" || req.Nickname == "" || req.Mobile == "" {
		return nil, status.Error(codes.FailedPrecondition, "")
	}
	// 加盐
	salt, encodedPwd := password.Encode(req.PassWord, global.Options)
	req.PassWord = fmt.Sprintf("$%s$%s$%s", global.HashMethodName, salt, encodedPwd)

	user := &userpb.UserInfo{
		NickName: req.Nickname,
		PassWord: req.PassWord,
		Mobile:   req.Mobile,
	}
	err := u.Dao.CreateUser(ctx, user)

	if err != nil {
		return nil, status.Error(err.Code, err.Err.Error())
	}
	return user, nil
}

func (u *UserService) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*emptypb.Empty, error) {
	user := &userpb.UserInfo{
		Id:       req.Id,
		NickName: req.NickName,
		Gender:   req.Gender,
		Birthday: req.Birthday,
	}
	if err := u.Dao.UpdateUser(ctx, user); err != nil {
		return nil, status.Error(err.Code, err.Err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (u *UserService) CheckPassWord(ctx context.Context, req *userpb.CheckPassWordRequest) (*userpb.CheckPassWordResponse, error) {
	pwd := strings.Split(req.EncPwd, "$")
	check := password.Verify(req.PassWord, pwd[2], pwd[3], global.Options)
	return &userpb.CheckPassWordResponse{Success: check}, nil
}

package service

import (
	"context"
	"go-micro-toDoList/pkg/errmsg"

	"go-micro-toDoList/user/internal/repository/dao"
	"go-micro-toDoList/user/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s UserService) UserLogin(ctx context.Context, req *pb.UserRequest) (*pb.UserDetailResponse, error) {
	resp := &pb.UserDetailResponse{}
	resp.Code = errmsg.SUCCESS
	user, err := dao.NewUserDao(ctx).GetUserInfo(req)
	if err != nil {
		return nil, err
	}

	resp.UserDetail = &pb.UserResponse{
		UserId:   user.UserID,
		NickName: user.NickName,
		UserName: user.UserName,
	}

	return resp, nil
}

package service

import (
	"context"

	"go-micro-toDoList/app/user/internal/repository/dao"
	"go-micro-toDoList/app/user/pb"
	"go-micro-toDoList/pkg/errmsg"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) UserLogin(ctx context.Context, req *pb.UserRequest) (*pb.UserDetailResponse, error) {
	resp := &pb.UserDetailResponse{}

	user, err := dao.NewUserDao(ctx).GetUserInfo(req)
	if err != nil {
		resp.Code = errmsg.FAILURE
		return resp, err
	}

	resp.UserDetail = &pb.UserResponse{
		UserId:   user.UserID,
		NickName: user.NickName,
		UserName: user.UserName,
	}

	resp.Code = errmsg.SUCCESS
	return resp, nil
}

func (s *UserService) UserRegister(ctx context.Context, req *pb.UserRequest) (*pb.UserCommonResponse, error) {
	resp := &pb.UserCommonResponse{}
	err := dao.NewUserDao(ctx).CreateUser(req)
	if err != nil {
		resp.Code = errmsg.FAILURE
		return resp, err
	}

	err = dao.NewUserDao(ctx).CreateUser(req)
	if err != nil {
		resp.Code = errmsg.FAILURE
		return resp, err
	}

	resp.Code = errmsg.SUCCESS
	resp.Data = errmsg.GetMsg(errmsg.SUCCESS)

	return resp, nil
}

func (s *UserService) UserLogout(ctx context.Context, req *pb.UserRequest) (resp *pb.UserCommonResponse, err error) {
	return
}
package service

import (
	"context"
	"micro-toDoList/app/user/internal/repository/dao"
	"micro-toDoList/pkg/errmsg"
	"micro-toDoList/pkg/pb/user_pb"
)

type UserService struct {
	user_pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) UserLogin(ctx context.Context, req *user_pb.UserRequest) (*user_pb.UserDetailResponse, error) {
	resp := &user_pb.UserDetailResponse{}

	// if user exists
	user, err := dao.NewUserDao(ctx).GetUserInfo(req)
	if err != nil {
		resp.Code = errmsg.FAILURE
		return resp, err
	}

	// if pwd correct
	if !user.CheckPwd(req.Password) {
		resp.Code = errmsg.FAILURE
		return resp, nil
	}

	resp.UserDetail = &user_pb.UserResponse{
		UserId:   user.UserID,
		NickName: user.NickName,
		UserName: user.UserName,
	}

	resp.Code = errmsg.SUCCESS
	return resp, nil
}

func (s *UserService) UserRegister(ctx context.Context, req *user_pb.UserRequest) (*user_pb.UserCommonResponse, error) {
	resp := &user_pb.UserCommonResponse{}
	err := dao.NewUserDao(ctx).CreateUser(req)
	if err != nil {
		resp.Code = errmsg.FAILURE
		return resp, err
	}

	// 学习这种写法；resp.msg不手动设置，返回去是不会有这一个field的
	resp.Code = errmsg.SUCCESS
	resp.Data = errmsg.GetMsg(errmsg.SUCCESS)

	return resp, nil
}

func (s *UserService) UserLogout(ctx context.Context, req *user_pb.UserRequest) (resp *user_pb.UserCommonResponse, err error) {
	return
}

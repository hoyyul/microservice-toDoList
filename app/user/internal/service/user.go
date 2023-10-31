package service

import (
	"context"
	"micro-toDoList/app/user/internal/repository/dao"
	"micro-toDoList/pkg/errmsg"
	"micro-toDoList/pkg/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) UserLogin(ctx context.Context, req *pb.UserRequest) (*pb.UserDetailResponse, error) {
	resp := &pb.UserDetailResponse{}

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

	resp.Code = errmsg.SUCCESS
	resp.Data = errmsg.GetMsg(errmsg.SUCCESS)

	return resp, nil
}

func (s *UserService) UserLogout(ctx context.Context, req *pb.UserRequest) (resp *pb.UserCommonResponse, err error) {
	return
}

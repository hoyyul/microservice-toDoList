package rpc

import (
	"context"
	"errors"
	"go-micro-toDoList/pkg/errmsg"
	"go-micro-toDoList/pkg/pb"
)

func UserLogin(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	// 这里接收的是detail res，但最终返回的是普通res
	r, err := UserClient.UserLogin(ctx, req)
	if err != nil {
		return nil, err
	}

	// pwd incorrect
	if r.Code == errmsg.FAILURE {
		err := errors.New("password not matched")
		return nil, err
	}

	return r.UserDetail, nil
}

func UserRegister(ctx context.Context, req *pb.UserRequest) (*pb.UserCommonResponse, error) {
	r, err := UserClient.UserRegister(ctx, req)
	if err != nil {
		return nil, err
	}

	return r, nil
}
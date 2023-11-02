package rpc

import (
	"context"
	"errors"
	"micro-toDoList/pkg/errmsg"
	"micro-toDoList/pkg/pb/user_pb"
)

func UserLogin(ctx context.Context, req *user_pb.UserRequest) (*user_pb.UserResponse, error) {
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

func UserRegister(ctx context.Context, req *user_pb.UserRequest) (*user_pb.UserCommonResponse, error) {
	r, err := UserClient.UserRegister(ctx, req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

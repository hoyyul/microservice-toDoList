package rpc

import (
	"context"
	"micro-toDoList/pkg/pb/task_pb"
)

func TaskCreate(ctx context.Context, req *task_pb.TaskRequest) (*task_pb.TaskCommonResponse, error) {
	r, err := TaskClient.TaskCreate(ctx, req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func TaskDelete(ctx context.Context, req *task_pb.TaskRequest) (*task_pb.TaskCommonResponse, error) {
	r, err := TaskClient.TaskDelete(ctx, req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func TaskUpdate(ctx context.Context, req *task_pb.TaskRequest) (*task_pb.TaskCommonResponse, error) {
	r, err := TaskClient.TaskUpdate(ctx, req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func TaskShow(ctx context.Context, req *task_pb.TaskRequest) ([]*task_pb.Task, error) {
	r, err := TaskClient.TaskShow(ctx, req)
	if err != nil {
		return nil, err
	}

	return r.GetTaskDetail(), nil
}

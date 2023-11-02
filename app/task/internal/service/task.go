package service

import (
	"context"
	"micro-toDoList/app/task/internal/repository/dao"
	"micro-toDoList/pkg/errmsg"
	"micro-toDoList/pkg/pb/task_pb"
)

type TaskService struct {
	task_pb.UnimplementedTaskServiceServer
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (s *TaskService) TaskCreate(ctx context.Context, req *task_pb.TaskRequest) (*task_pb.TaskCommonResponse, error) {
	resp := &task_pb.TaskCommonResponse{}
	err := dao.NewTaskDao(ctx).CreateTask(req)
	if err != nil {
		resp.Code = errmsg.FAILURE
		return nil, err
	}

	resp.Code = errmsg.SUCCESS
	resp.Data = errmsg.GetMsg(errmsg.SUCCESS)
	return resp, nil
}

func (s *TaskService) TaskDelete(ctx context.Context, req *task_pb.TaskRequest) (*task_pb.TaskCommonResponse, error) {
	resp := &task_pb.TaskCommonResponse{}
	err := dao.NewTaskDao(ctx).DeleteTaskById(req.TaskId)
	if err != nil {
		resp.Code = errmsg.FAILURE
		return nil, err
	}

	resp.Code = errmsg.SUCCESS
	resp.Data = errmsg.GetMsg(errmsg.SUCCESS)
	return resp, nil
}

func (s *TaskService) TaskUpdate(ctx context.Context, req *task_pb.TaskRequest) (*task_pb.TaskCommonResponse, error) {
	resp := &task_pb.TaskCommonResponse{}
	err := dao.NewTaskDao(ctx).UpdateTask(req)
	if err != nil {
		resp.Code = errmsg.FAILURE
		return nil, err
	}

	resp.Code = errmsg.SUCCESS
	resp.Data = errmsg.GetMsg(errmsg.SUCCESS)
	return resp, nil
}

func (s *TaskService) TaskShow(ctx context.Context, req *task_pb.TaskRequest) (*task_pb.TaskDetailResponse, error) {
	resp := &task_pb.TaskDetailResponse{}
	tasks, err := dao.NewTaskDao(ctx).GetTaskByUserId(req.UserId)
	if err != nil {
		resp.Code = errmsg.FAILURE
		return nil, err
	}

	resp.Code = errmsg.SUCCESS

	for i := range tasks {
		resp.TaskDetail = append(resp.TaskDetail, &task_pb.Task{
			TaskId:    tasks[i].TaskId,
			UserId:    tasks[i].UserId,
			Status:    int64(tasks[i].Status),
			Title:     tasks[i].Title,
			Content:   tasks[i].Content,
			StartTime: tasks[i].StartTime,
			EndTime:   tasks[i].EndTime,
		})
	}

	return resp, nil
}

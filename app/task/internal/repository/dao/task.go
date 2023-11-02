package dao

import (
	"context"
	"errors"
	"micro-toDoList/app/task/internal/repository/model"
	"micro-toDoList/pkg/pb/task_pb"

	"gorm.io/gorm"
)

type TaskDao struct {
	*gorm.DB
}

func NewTaskDao(ctx context.Context) *TaskDao {
	return &TaskDao{DBWithContext(ctx)}
}

func (dao *TaskDao) CreateTask(req *task_pb.TaskRequest) error {
	var count int64
	dao.Model(model.Task{}).Where("task_title=?", req.Title).Count(&count)
	if count > 0 {
		return errors.New("task already exists")
	}
	task := &model.Task{
		UserId:    req.UserId,
		Status:    int(req.Status),
		Title:     req.Title,
		Content:   req.Content,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	err := dao.Model(model.Task{}).Create(&task).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *TaskDao) DeleteTaskById(taskId int64) error {
	err := dao.Model(model.Task{}).Where("task_id=?", taskId).Delete(model.Task{}).Error
	if err != nil {
		return err
	}
	return err
}

func (dao *TaskDao) UpdateTask(req *task_pb.TaskRequest) error {
	var count int64
	dao.Model(model.Task{}).Where("task_id=?", req.TaskId).Count(&count)
	if count == 0 {
		return errors.New("task doesn't exist")
	}

	taskMap := make(map[string]interface{})
	taskMap["title"] = req.Title
	taskMap["content"] = req.Content
	taskMap["status"] = int(req.Status)
	taskMap["start_time"] = req.StartTime
	taskMap["end_time"] = req.EndTime

	err := dao.Model(model.Task{}).Where("task_id=?", req.TaskId).Updates(&taskMap).Error
	if err != nil {
		return err
	}

	return nil
}

func (dao *TaskDao) GetTaskByUserId(userId int64) (tasks []*model.Task, err error) {
	err = dao.Model(model.Task{}).Where("user_id=?", userId).Find(&tasks).Error
	return
}

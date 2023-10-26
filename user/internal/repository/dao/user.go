package dao

import (
	"context"
	"go-micro-toDoList/user/internal/repository/model"
	"go-micro-toDoList/user/pb"

	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{DBWithContext(ctx)}
}

func (dao *UserDao) GetUserInfo(req *pb.UserRequest) (user *model.User, err error) {
	err = dao.Where("user_name=?", req.UserName).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return
}

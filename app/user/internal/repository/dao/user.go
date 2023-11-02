package dao

import (
	"context"
	"errors"
	"micro-toDoList/app/user/internal/repository/model"
	"micro-toDoList/pkg/pb/user_pb"

	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{DBWithContext(ctx)}
}

func (dao *UserDao) GetUserInfo(req *user_pb.UserRequest) (user *model.User, err error) {
	err = dao.Where("user_name=?", req.UserName).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return
}

// 1.创建user不用指定userid，因为gorm会自动创建
// 2.增删改前需要考虑user是否已经存在
func (dao *UserDao) CreateUser(req *user_pb.UserRequest) error {
	var count int64
	dao.Model(&model.User{}).Where("user_name=?", req.UserName).Count(&count)
	if count > 0 {
		return errors.New("user already exists")
	}

	user := &model.User{
		UserName: req.UserName,
		NickName: req.NickName,
	}
	err := user.SetPwd(req.Password)
	if err != nil {
		return err
	}

	err = dao.Model(&model.User{}).Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

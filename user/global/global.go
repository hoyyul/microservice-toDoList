package global

import (
	"go-micro-toDoList/user/config"

	"github.com/sirupsen/logrus"
)

var (
	Config *config.Config
	//Mysql  *gorm.DB
	Logger *logrus.Logger
)

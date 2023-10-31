package global

import (
	"micro-toDoList/config"

	"github.com/sirupsen/logrus"
)

var (
	Config *config.Config
	Logger *logrus.Logger
)

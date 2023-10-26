package main

import (
	"go-micro-toDoList/user/internal/repository/dao"
	"go-micro-toDoList/user/setting"
)

func main() {
	setting.InitConfig()
	setting.InitLogger()
	dao.InitDB()
}

package main

import (
	"go-micro-toDoList/app/user/internal/repository/dao"
	"go-micro-toDoList/app/user/internal/service"

	"go-micro-toDoList/global"
	"go-micro-toDoList/pkg/eTcd"
	"go-micro-toDoList/pkg/pb"
	"go-micro-toDoList/setting"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// global initial setting
	setting.InitConfig()
	setting.InitLogger()

	// local initial setting
	dao.InitDB()

	// register service to etcd
	etcdRegistrar := eTcd.NewRegistrar(global.Config.Etcd.Address) // etcd
	defer etcdRegistrar.UnRegister()

	server := eTcd.Server{
		Name: global.Config.Services["user"].Name,
		Addr: global.Config.Services["user"].Address,
	}
	if err := etcdRegistrar.Register(server, 10); err != nil { // my server
		global.Logger.Fatalln("Failed to register service to Etcd Server")
		return
	}

	// bind grpc service
	listener, err := net.Listen("tcp", global.Config.Server.Addr)
	if err != nil {
		return
	}
	global.Logger.Println("Listeing to %v", listener.Addr())

	grpcSrv := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcSrv, service.NewUserService())

	// start grpc server
	if err := grpcSrv.Serve(listener); err != nil {
		return
	}
	global.Logger.Println("User server is running on %s", global.Config.Server.Addr)
}

package main

import (
	"go-micro-toDoList/pkg/eTcd"
	"go-micro-toDoList/user/global"
	"go-micro-toDoList/user/internal/service"
	"go-micro-toDoList/user/pb"
	"go-micro-toDoList/user/setting"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// initial setting
	setting.InitConfig()
	setting.InitLogger()

	// register service to etcd
	etcdRegistrar := eTcd.NewRegistrar(global.Config.Server.Addr)
	defer etcdRegistrar.UnRegister()

	taskNode := eTcd.Server{
		Name: global.Config.Server.Name,
		Addr: global.Config.Server.Addr,
	}
	if err := etcdRegistrar.Register(taskNode, 10); err != nil {
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

	// start to serve
	if err := grpcSrv.Serve(listener); err != nil {
		return
	}
	global.Logger.Println("User server is running on %s", global.Config.Server.Addr)
}

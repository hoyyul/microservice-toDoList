package main

import (
	"micro-toDoList/app/user/internal/repository/dao"
	"micro-toDoList/app/user/internal/service"
	"micro-toDoList/global"
	"micro-toDoList/pkg/eTcd"
	"micro-toDoList/pkg/pb"
	"micro-toDoList/setting"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// global initial setting
	setting.InitConfig()
	setting.InitLogger()

	// local initial setting
	dao.InitDB()

	userName := global.Config.Services["user"].Name
	userAddr := global.Config.Services["user"].Address

	// register service to etcd
	etcdRegistrar := eTcd.NewRegistrar(global.Config.Etcd.Address) // etcd
	defer etcdRegistrar.UnRegister()

	server := eTcd.Server{
		Name: userName,
		Addr: userAddr,
	}
	if err := etcdRegistrar.Register(server, 10); err != nil { // my server
		global.Logger.Fatalln("Failed to register service to Etcd Server")
		return
	}

	// bind grpc service
	listener, err := net.Listen("tcp", userAddr)
	if err != nil {
		global.Logger.Panic(err)
		return
	}
	global.Logger.Printf("Listeing to %v", listener.Addr())

	grpcSrv := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcSrv, service.NewUserService())

	// start grpc server
	if err := grpcSrv.Serve(listener); err != nil {
		return
	}
	global.Logger.Printf("User server is running on %s", userAddr)

}

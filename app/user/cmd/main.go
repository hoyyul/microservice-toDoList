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

	// local server initial setting
	dao.InitDB()
	userName := global.Config.Services["user"].Name
	userAddr := global.Config.Services["user"].Address
	etcdAddrs := []string{global.Config.Etcd.Address}
	etcdRegistrar := eTcd.NewRegistrar(etcdAddrs)
	defer etcdRegistrar.Stop()
	server := eTcd.Server{
		Name: userName,
		Addr: userAddr,
	}

	// register service to etcd
	if err := etcdRegistrar.Register(server, 10); err != nil {
		global.Logger.Fatalln("Failed to register service to Etcd Server")
		return
	}

	// listen to server
	listener, err := net.Listen("tcp", userAddr)
	if err != nil {
		global.Logger.Panic(err)
		return
	}
	global.Logger.Printf("Listeing to %v", listener.Addr())

	// register my server to grpc server
	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()
	pb.RegisterUserServiceServer(grpcServer, service.NewUserService())

	// enable grpc  to serve
	if err := grpcServer.Serve(listener); err != nil {
		return
	}
	global.Logger.Printf("User server is running on %s", userAddr)

}

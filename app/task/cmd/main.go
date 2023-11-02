package main

import (
	"micro-toDoList/app/task/internal/repository/dao"
	"micro-toDoList/app/task/internal/service"
	"micro-toDoList/global"
	"micro-toDoList/pkg/eTcd"
	"micro-toDoList/pkg/pb/task_pb"
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
	taskName := global.Config.Services["task"].Name
	taskAddr := global.Config.Services["task"].Address
	etcdAddrs := []string{global.Config.Etcd.Address}
	etcdRegistrar := eTcd.NewRegistrar(etcdAddrs)
	defer etcdRegistrar.Stop()
	server := eTcd.Server{
		Name: taskName,
		Addr: taskAddr,
	}

	// register task server address to etcd
	if err := etcdRegistrar.Register(server, 10); err != nil {
		global.Logger.Fatalln("Failed to register service to Etcd Server")
		return
	}

	// listen to server
	listener, err := net.Listen("tcp", taskAddr)
	if err != nil {
		global.Logger.Panic(err)
		return
	}
	defer listener.Close()
	global.Logger.Printf("Listeing to %v", listener.Addr())

	// register my server to grpc server
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()
	task_pb.RegisterTaskServiceServer(grpcServer, service.NewTaskService())

	// enable grpc to serve
	if err := grpcServer.Serve(listener); err != nil {
		return
	}
	global.Logger.Printf("User server is running on %s", taskAddr)

}

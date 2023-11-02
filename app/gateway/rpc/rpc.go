package rpc

import (
	"context"
	"fmt"
	"micro-toDoList/global"
	"micro-toDoList/pkg/eTcd"
	"micro-toDoList/pkg/pb/user_pb"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

var (
	// 学习一下上下文
	ctx    context.Context
	cancel context.CancelFunc

	Resolver *eTcd.Resolver

	UserClient user_pb.UserServiceClient
)

func Init() {
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// grpc register a etcd resolver
	etcdAddrs := []string{global.Config.Etcd.Address}
	Resolver = eTcd.NewResolver(etcdAddrs)
	resolver.Register(Resolver)

	// connect etcd and start discovery service
	initGrpcClient(global.Config.Services["user"].Name, &UserClient)
}

// create a grpc client sub for a connection
func initGrpcClient(service string, client interface{}) {
	conn, err := enableEtcdDiscovery(service)
	if err != nil {
		global.Logger.Panicln(err)
	}

	switch c := client.(type) { //记住这种写发，c := client.(type)实际判断是client（指针）的类型；*c = pb.NewUserServiceClient(conn)实际是把值赋到client指向的地址。也就是说赋给UserClient
	case *user_pb.UserServiceClient:
		*c = user_pb.NewUserServiceClient(conn)
	default:
		global.Logger.Panicln("Invalid type")
	}
}

func enableEtcdDiscovery(serviceName string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	target := fmt.Sprintf("%s:///%s", Resolver.Scheme(), serviceName)

	//  enable loadbalance policy
	if global.Config.Services[serviceName].LoadBalance {
		global.Logger.Printf("Enable load balance for %s service", serviceName)
		opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	}

	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

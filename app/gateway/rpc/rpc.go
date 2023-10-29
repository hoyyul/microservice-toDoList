package rpc

import (
	"context"
	"fmt"
	"go-micro-toDoList/global"
	"go-micro-toDoList/pkg/pb"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	resolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// 学习一下上下文
	ctx    context.Context
	cancel context.CancelFunc

	UserClient pb.UserServiceClient
)

func Init() {
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)

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
	case *pb.UserServiceClient:
		*c = pb.NewUserServiceClient(conn)
	default:
		global.Logger.Panicln("Invalid type")
	}
}

func enableEtcdDiscovery(service string) (*grpc.ClientConn, error) {
	// create a etcd client
	etcdUrl := fmt.Sprintf("http://%s", global.Config.Etcd.Address)
	etcdClient, _ := clientv3.NewFromURL(etcdUrl)

	// create a ectd resolver
	etcdResolver, _ := resolver.NewBuilder(etcdClient)

	// connect
	opts := []grpc.DialOption{ // client-side option; serviceOption is server-side option
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(etcdResolver),
	}

	//  enable loadbalance policy
	if global.Config.Services[service].LoadBalance {
		global.Logger.Println("Enable load balance for %s service", service)
		opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	}

	// discovery service
	addr := fmt.Sprintf("etcd:///%s/", service) // 在client和server间加入中间键来交流
	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

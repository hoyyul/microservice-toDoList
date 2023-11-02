package eTcd

import (
	"context"
	"micro-toDoList/global"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

var myScheme string = "etcd"

type Resolver struct {
	EtcdAddrs    []string
	ServiceAddrs []resolver.Address
	Cli          *clientv3.Client
	KeyPrefix    string
	DialTimeout  int
	ClientConn   resolver.ClientConn // 可以理解成grpc客户端和etcd resolver的交互

	watchCh clientv3.WatchChan
	closeCh chan struct{}
}

func NewResolver(etcdAddrs []string) *Resolver {
	return &Resolver{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 5,
		closeCh:     make(chan struct{}),
	}
}

// impletement resolver.Revolver interface
// 1. ResolveNow
// 2. Close
func (r *Resolver) ResolveNow(op resolver.ResolveNowOptions) {}

func (r *Resolver) Close() {
	r.closeCh <- struct{}{}
}

// impletement resolver.Builder interface
// 1. Build
// 2. Scheme
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// connect to etcd server
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	r.Cli = cli

	r.ClientConn = cc

	// extract key prefix from gprc target; 其实target 是来自于 grpc.dial(target)
	r.KeyPrefix = BuildPrefix(Server{Name: target.Endpoint()}) // e.g user -> /user/

	// get service addresses
	if err := r.synchro(); err != nil {
		return r, err
	}

	// keep watching
	go r.watch()

	return r, nil
}

func (r *Resolver) Scheme() string {
	return myScheme
}

// watching service key-value update
func (r *Resolver) watch() {
	ticker := time.NewTicker(time.Second * time.Duration(60))
	r.watchCh = r.Cli.Watch(context.Background(), r.KeyPrefix, clientv3.WithPrefix()) // watchCh : <-watchRepsonse

	for {
		select {
		// end discovery
		case <-r.closeCh:
			// close connect
			return
		// update event occurs
		case resp, ok := <-r.watchCh:
			if ok {
				err := r.update(resp.Events)
				if err != nil {
					global.Logger.Error(err)
				}
			}
		// update service discovery address very min
		case <-ticker.C:
			if err := r.synchro(); err != nil {
				global.Logger.Error(err)
			}
		}
	}
}

// update action
func (r *Resolver) update(events []*clientv3.Event) error {
	for _, event := range events {
		srvInfo, err := ParseValue(event.Kv.Value)
		if err != nil {
			return err
		}
		addr := resolver.Address{Addr: srvInfo.Addr}

		switch event.Type {
		case clientv3.EventTypePut:
			// in this app; only having add new service logic
			if !exist(addr, r.ServiceAddrs) {
				r.ServiceAddrs = append(r.ServiceAddrs, addr)
			}
			// todo: update logic...
		case clientv3.EventTypeDelete:
			if temp, ok := remove(r.ServiceAddrs, addr); ok {
				r.ServiceAddrs = temp
			}
		}
		r.ClientConn.UpdateState(resolver.State{Addresses: r.ServiceAddrs})
	}
	return nil
}

// update to new service addresses
func (r *Resolver) synchro() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	// clean old addresses
	r.ServiceAddrs = []resolver.Address{}

	// get service key-value pairs by key prefix
	resp, err := r.Cli.Get(ctx, r.KeyPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	// take the addr stored in value
	for _, kv := range resp.Kvs {
		srvInfo, err := ParseValue(kv.Value)
		if err != nil {
			return err
		}
		// add to service addresslist
		r.ServiceAddrs = append(r.ServiceAddrs, resolver.Address{Addr: srvInfo.Addr})
	}

	// update the new discovery service list to etcd
	err = r.ClientConn.UpdateState(resolver.State{Addresses: r.ServiceAddrs})
	if err != nil {
		return err
	}

	return nil
}

// helper func
func exist(addr resolver.Address, addrs []resolver.Address) bool {
	for i := range addrs {
		if addrs[i] == addr {
			return true
		}
	}
	return false
}

func remove(s []resolver.Address, addr resolver.Address) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr.Addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}

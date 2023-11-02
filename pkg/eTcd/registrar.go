package eTcd

import (
	"context"
	"encoding/json"
	"fmt"
	"micro-toDoList/global"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Registrar struct {
	EtcdAddrs   []string
	DialTimeout int

	Cli    *clientv3.Client
	Srv    Server
	SrvTTL int64

	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse
	closeCh     chan struct{}
	LeaseID     clientv3.LeaseID
}

func NewRegistrar(etcdAddrs []string) *Registrar {
	return &Registrar{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 5,
		closeCh:     make(chan struct{}),
	}
}

func (r *Registrar) Register(srv Server, ttl int64) error {
	// save server info
	r.Srv = srv

	// create an etcd client to connect etcd server
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		global.Logger.Panicln("Failed to connect etcd server")
		return err
	}
	r.Cli = cli

	// register service
	go r.keepAlive()

	return nil
}

func (r *Registrar) register() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	// 1. Grant; create lease
	lease, err := r.Cli.Grant(ctx, r.SrvTTL)
	if err != nil {
		return err
	}

	// 2. KeepAlive; update the lease for a time period
	r.LeaseID = lease.ID
	r.keepAliveCh, err = r.Cli.KeepAlive(ctx, r.LeaseID)
	if err != nil {
		return err
	}

	// 3. Put; save key-value in etcd, etcd is a key-value store
	key := BuildRegisterPath(r.Srv) // key; e.g /user/127.0.0.1:10001
	fmt.Println(key)
	srvInfo, err := json.Marshal(r.Srv) // value
	if err != nil {
		return err
	}

	_, err = r.Cli.Put(context.Background(), key, string(srvInfo), clientv3.WithLease(r.LeaseID))
	if err != nil {
		return err
	}

	return err
}

func (r *Registrar) keepAlive() {
	ticker := time.NewTicker(time.Duration(r.DialTimeout) * time.Second)
	for {
		select { // channel 版本的switch；看看那个准备好了； 如果多个准备好了就随机选一个
		// end user service
		case <-r.closeCh:
			if err := r.unRegister(); err != nil {
				global.Logger.Error(err)
			}
		// release the capacity of keepAlive channel
		case resp := <-r.keepAliveCh:
			if resp == nil {
				if err := r.register(); err != nil {
					global.Logger.Error(err)
				}
			}
		// periodacally check service key-value alive in etcd
		case <-ticker.C:
			resp, err := r.Cli.Get(context.Background(), BuildRegisterPath(r.Srv))
			if err != nil {
				if err := r.register(); err != nil {
					global.Logger.Error(err)
				}
			}
			if resp.Count == 0 {
				if err := r.register(); err != nil {
					global.Logger.Error(err)
				}
			}
		}

	}

}

func (r *Registrar) Stop() {
	r.closeCh <- struct{}{}
}

func (r *Registrar) unRegister() error {
	// delete registered service
	_, err := r.Cli.Delete(context.Background(), BuildRegisterPath(r.Srv))
	if err != nil {
		return err
	}

	// revoke lease
	_, err = r.Cli.Revoke(context.Background(), r.LeaseID)
	if err != nil {
		return err
	}

	// close connect
	r.Cli.Close()

	return nil
}

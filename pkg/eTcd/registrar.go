package eTcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type Registrar struct {
	EtcdAddr string

	// hiden info
	cli *clientv3.Client
	srv Server
}

func NewRegistrar(etcdAddr string) *Registrar {
	return &Registrar{EtcdAddr: etcdAddr}
}

func (r *Registrar) Register(srv Server, ttl int64) error {
	// set server
	r.srv = srv

	// create etcd conn
	cli, err := clientv3.NewFromURL(fmt.Sprintf("http://%s", r.EtcdAddr))
	if err != nil {
		return err
	}
	r.cli = cli

	// grant lease
	lease, err := r.cli.Grant(context.TODO(), ttl)
	if err != nil {
		return err
	}

	// add service to endpoint
	em, err := endpoints.NewManager(r.cli, BuildPrefix(r.srv))
	if err != nil {
		return err
	}
	em.AddEndpoint(context.TODO(), BuildRegisterPath(r.srv), endpoints.Endpoint{Addr: r.srv.Addr}, clientv3.WithLease(lease.ID))

	// revoke
	ch, _ := r.cli.KeepAlive(context.TODO(), lease.ID)

	// consume the ka in channel to avoid max capacity
	go func() {
		for {
			ka := <-ch
			if ka != nil {
				continue
			}
		}
	}()

	return nil
}

func (r *Registrar) UnRegister() error {
	em, err := endpoints.NewManager(r.cli, r.srv.Name)
	if err != nil {
		return err
	}

	return em.DeleteEndpoint(context.TODO(), fmt.Sprintf("%s:%s", r.srv.Name, r.srv.Addr))
}

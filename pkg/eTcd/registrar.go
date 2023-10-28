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
	cli, err := clientv3.NewFromURL(r.EtcdAddr)
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
	em, err := endpoints.NewManager(r.cli, r.srv.Name)
	if err != nil {
		return err
	}
	em.AddEndpoint(context.TODO(), fmt.Sprintf("%s:%s", r.srv.Name, r.srv.Addr), endpoints.Endpoint{Addr: r.srv.Addr}, clientv3.WithLease(lease.ID))

	// revoke
	go r.cli.KeepAlive(context.TODO(), lease.ID)

	return nil
}

func (r *Registrar) UnRegister() error {
	em, err := endpoints.NewManager(r.cli, r.srv.Name)
	if err != nil {
		return err
	}

	return em.DeleteEndpoint(context.TODO(), fmt.Sprintf("%s:%s", r.srv.Name, r.srv.Addr))
}

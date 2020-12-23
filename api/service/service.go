package service

import (
	"context"
	"github.com/nevercase/k8s-controller-custom-resource/api/conf"
	"github.com/nevercase/k8s-controller-custom-resource/api/group"
)

type Service interface {
	Listen()
	Close()
}

type service struct {
	conf conf.Config

	conn   ConnHub
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *service) Listen() {
	go s.initWSService(s.conf.ApiService())
}

func (s *service) Close() {
	s.cancel()
}

func NewService(c conf.Config) Service {
	ctx, cancel := context.WithCancel(context.Background())
	g := group.NewGroup(ctx, c.MasterUrl(), c.KubeConfig(), c.DockerHub())
	s := &service{
		conf:   c,
		conn:   NewConnHub(ctx, g),
		ctx:    ctx,
		cancel: cancel,
	}
	return s
}

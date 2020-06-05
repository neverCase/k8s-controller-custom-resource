package service

import (
	"context"

	"github.com/nevercase/k8s-controller-custom-resource/api/conf"
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
	s := &service{
		conf:   c,
		conn:   NewConnHub(ctx),
		ctx:    ctx,
		cancel: cancel,
	}
	return s
}

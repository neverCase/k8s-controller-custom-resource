package service

import (
	"context"
)

type Service interface {
	Start()
	Close()
}

type service struct {
	conn   ConnHub
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *service) Start() {
	go s.initWSService("0.0.0.0:9090")
}

func (s *service) Close() {
	s.cancel()
}

func NewService() Service {
	ctx, cancel := context.WithCancel(context.Background())
	s := &service{
		conn:   NewConnHub(ctx),
		ctx:    ctx,
		cancel: cancel,
	}
	return s
}

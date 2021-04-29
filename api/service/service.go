package service

import (
	"context"
	"github.com/Shanghai-Lunara/pkg/authentication"
	"github.com/Shanghai-Lunara/pkg/casbinrbac"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nevercase/k8s-controller-custom-resource/api/conf"
	"github.com/nevercase/k8s-controller-custom-resource/api/group"
	"net/http"
)

type Service interface {
	Close()
}

type service struct {
	conf   conf.Config
	conn   ConnHub
	server *http.Server
	ctx    context.Context
	cancel context.CancelFunc
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
	router := gin.New()
	router.Use(cors.Default())
	casbinrbac.NewWithMysqlConf(c.RbacRulePath(), c.RbacMysqlPath(), "/rbac", router)
	authentication.New("/authentication", router)
	router.Group("dashboard").GET("", s.handler)
	server := &http.Server{
		Addr:    s.conf.ApiService(),
		Handler: router,
	}
	s.server = server
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				zaplogger.Sugar().Info("Server closed under request")
			} else {
				zaplogger.Sugar().Info("Server closed unexpected err:", err)
			}
		}
	}()
	return s
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *service) handler(c *gin.Context) {
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	s.conn.NewConn(conn)
}

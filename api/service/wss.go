package service

import (
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"k8s.io/klog"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *service) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		klog.V(2).Info(err)
		return
	}
	s.conn.NewConn(conn)
}

func (s *service) initWSService(addr string) {
	klog.Info("initWSService")
	http.HandleFunc("/", s.wsHandler)
	l, err := net.Listen("tcp4", addr)
	if err != nil {
		klog.Fatal(err)
	}
	klog.Fatal(http.Serve(l, nil))
}

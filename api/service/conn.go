package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	//apiV1 "k8s.io/api/core/v1"
	//metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog"

	"github.com/nevercase/k8s-controller-custom-resource/api/group"
	"github.com/nevercase/k8s-controller-custom-resource/api/proto"
	mysqlOperatorV1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
)

type ConnHub interface {
	NewConn(conn *websocket.Conn)
}

type connHub struct {
	group group.Group

	mu           sync.RWMutex
	autoClientId int32
	connections  map[int32]Conn
	ctx          context.Context
}

func (ch *connHub) NewConn(conn *websocket.Conn) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	id := atomic.AddInt32(&ch.autoClientId, 1)
	ch.connections[id] = NewConn(id, ch.ctx, conn, ch.group)
	go func() {
		if err := ch.connections[id].ReadPump(); err != nil {
			klog.V(2).Info(err)
		}
	}()
	go func() {
		if err := ch.connections[id].WritePump(); err != nil {
			klog.V(2).Info(err)
		}
	}()
	go ch.connections[id].KeepAlive()
}

func NewConnHub(ctx context.Context, g group.Group) ConnHub {
	return &connHub{
		group:        g,
		autoClientId: 0,
		connections:  make(map[int32]Conn, 0),
		ctx:          ctx,
	}
}

type Conn interface {
	Ping()
	KeepAlive()
	ReadPump() (err error)
	SendToChannel(data interface{}) (err error)
	WritePump() (err error)
	Close()
}

func NewConn(clientId int32, ctx context.Context, ws *websocket.Conn, g group.Group) Conn {
	c := &conn{
		group:             g,
		clientId:          clientId,
		conn:              ws,
		readChan:          make(chan interface{}),
		writeChan:         make(chan interface{}),
		lastHeartBeatTime: time.Now(),
		status:            connAlive,
		ctx:               ctx,
	}
	return c
}

type connStatus int

const (
	connAlive  connStatus = iota
	connClosed connStatus = 1
)

type conn struct {
	group group.Group

	mu                sync.RWMutex
	clientId          int32
	conn              *websocket.Conn
	readChan          chan interface{}
	writeChan         chan interface{}
	lastHeartBeatTime time.Time
	tick              *time.Ticker
	status            connStatus
	ctx               context.Context
	cancel            context.CancelFunc
}

func (c *conn) Ping() {
	c.lastHeartBeatTime = time.Now()
}

func (c *conn) KeepAlive() {
	defer c.Close()
	for {
		tick := time.NewTicker(10 * time.Second)
		select {
		case <-c.ctx.Done():
			return
		case <-tick.C:
			if time.Now().Sub(c.lastHeartBeatTime) > 30 {
				klog.Info("keepAlive timeout")
				return
			}
		}
	}
}

func (c *conn) ReadPump() (err error) {
	defer c.Close()
	for {
		var msg proto.Request
		messageType, message, err := c.conn.ReadMessage()
		klog.Infof("messageType: %d message: %s err: %s\n", messageType, message, err)
		if err != nil {
			klog.V(2).Info(err)
			return err
		}
		if err = json.Unmarshal(message, &msg); err != nil {
			klog.V(2).Info(err)
			return err
		}
		klog.Info(msg.Service)

		//if t, err := c.group.Mysql().MysqloperatorV1().MysqlOperators(apiV1.NamespaceDefault).List(metaV1.ListOptions{}); err != nil {
		//	klog.V(2).Info(err)
		//} else {
		//	if err = c.SendToChannel(proto.GetResponse(t)); err != nil {
		//		klog.V(2.).Info(err)
		//		return err
		//	}
		//}

		klog.Info("1111111111")
		s := "svc-121"
		var a = &proto.List{Code: 1234, Result: s}
		if res, err := a.Marshal(); err != nil {
			klog.V(2).Info(err)
		} else {
			klog.Info("res List:", string(res))
			if err = c.SendToChannel(proto.GetResponse(string(res))); err != nil {
				return err
			}
		}

		var e = &proto.Mysql{Mysql: mysqlOperatorV1.MysqlOperator{}}
		if res, err := e.Marshal(); err != nil {
			klog.V(2).Info(err)
		} else {
			klog.Info("res mysql:", string(res))
			if err = c.SendToChannel(proto.GetResponse(string(res))); err != nil {
				return err
			}
		}

		switch msg.Service {
		case proto.SvcPing:
			c.Ping()
			if err = c.SendToChannel(proto.GetResponse(msg.Data)); err != nil {
				return err
			}
		case proto.SvcList:
			if err = c.SendToChannel(proto.GetResponse(msg.Data)); err != nil {
				return err
			}
		case proto.SvcWatch:
		case proto.SvcAdd:
		case proto.SvcUpdate:
		case proto.SvcDelete:
		}
	}
}

func (c *conn) SendToChannel(msg interface{}) (err error) {
	if c.status == connClosed {
		return
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.status == connClosed {
		return
	}
	tick := time.NewTicker(time.Second * 2)
	defer tick.Stop()
	for {
		select {
		case c.writeChan <- msg:
			return
		case <-tick.C:
			err = fmt.Errorf("wsSend timeout ws.cid:%d msg:(%v) ws:%v\n", c.clientId, msg, c)
			klog.Info(err)
			return err
		}
	}
}

func (c *conn) WritePump() (err error) {
	defer c.Close()
	for {
		select {
		case <-c.ctx.Done():
			return nil
		case msg, isClose := <-c.writeChan:
			if !isClose {
				return nil
			}
			log.Println("send to:", c.clientId, " msg:", msg)
			if err := c.conn.WriteJSON(msg); err != nil {
				return err
			}
		}
	}
}

func (c *conn) Close() {
	if c.status == connClosed {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status = connClosed
	if err := c.conn.Close(); err != nil {
		klog.V(2).Info(err)
	}
	close(c.writeChan)
	close(c.readChan)
}

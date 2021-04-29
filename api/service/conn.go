package service

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"k8s.io/klog/v2"

	"github.com/nevercase/k8s-controller-custom-resource/api/group"
	"github.com/nevercase/k8s-controller-custom-resource/api/handle"
	"github.com/nevercase/k8s-controller-custom-resource/api/proto"
)

type ConnHub interface {
	NewConn(conn *websocket.Conn)
}

type connHub struct {
	group  group.Group
	handle handle.Handle

	mu           sync.RWMutex
	autoClientId int32
	connections  map[int32]WsConn

	broadcast chan []byte

	ctx context.Context
}

func (ch *connHub) NewConn(conn *websocket.Conn) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	id := atomic.AddInt32(&ch.autoClientId, 1)
	ch.connections[id] = NewConn(id, ch.ctx, conn, ch.group, ch.handle)
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

func (ch *connHub) BroadcastWatch() {
	for {
		select {
		case msg := <-ch.broadcast:
			ch.mu.RLock()
			for _, v := range ch.connections {
				if err := v.SendToChannel(msg); err != nil {
					klog.V(2).Info("BroadcastWatch err:%v", err)
				}
			}
			ch.mu.RUnlock()
		}
	}
}

func NewConnHub(ctx context.Context, g group.Group) ConnHub {
	b := make(chan []byte, 4096)
	ch := &connHub{
		group:        g,
		handle:       handle.NewHandle(g, b),
		autoClientId: 0,
		connections:  make(map[int32]WsConn, 0),
		broadcast:    b,
		ctx:          ctx,
	}
	go ch.BroadcastWatch()
	return ch
}

type WsConn interface {
	Ping()
	KeepAlive()
	ReadPump() (err error)
	SendToChannel(data []byte) (err error)
	WritePump() (err error)
	Close()
}

func NewConn(clientId int32, ctx context.Context, ws *websocket.Conn, g group.Group, h handle.Handle) WsConn {
	c := &wsConn{
		handle:            h,
		group:             g,
		clientId:          clientId,
		conn:              ws,
		readChan:          make(chan interface{}),
		writeChan:         make(chan []byte),
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

type wsConn struct {
	group  group.Group
	handle handle.Handle

	mu                sync.RWMutex
	clientId          int32
	conn              *websocket.Conn
	readChan          chan interface{}
	writeChan         chan []byte
	lastHeartBeatTime time.Time
	tick              *time.Ticker
	status            connStatus
	ctx               context.Context
	cancel            context.CancelFunc
}

func (c *wsConn) Ping() {
	c.lastHeartBeatTime = time.Now()
}

func (c *wsConn) KeepAlive() {
	defer c.Close()
	for {
		tick := time.NewTicker(10 * time.Second)
		select {
		case <-c.ctx.Done():
			return
		case <-tick.C:
			if time.Now().Sub(c.lastHeartBeatTime) > 30*time.Second {
				klog.Info("keepAlive timeout")
				return
			}
		}
	}
}

func (c *wsConn) ReadPump() (err error) {
	defer c.Close()
	for {
		var msg proto.Request
		var res []byte
		_, message, err := c.conn.ReadMessage()
		//klog.Infof("messageType: %d message: %v err: %s\n", messageType, message, err)
		if err != nil {
			klog.V(2).Info(err)
			return err
		}
		if err = msg.Unmarshal(message); err != nil {
			klog.V(2).Info(err)
			return err
		}
		klog.Info("proto Request:", msg)
		switch proto.ApiService(msg.Param.Service) {
		case proto.SvcPing:
			c.Ping()
			if res, err = proto.GetResponse(msg.Param, []byte("ping success")); err != nil {
				return err
			}
		case proto.SvcCreate:
			if res, err = c.handle.KubernetesApi().Create(msg.Param, msg.Data); err != nil {
				klog.V(2).Info(err)
				if res, err = proto.ErrorResponse(msg.Param, err.Error()); err != nil {
					klog.V(2).Info(err)
				}
			}
		case proto.SvcUpdate:
			if res, err = c.handle.KubernetesApi().Update(msg.Param, msg.Data); err != nil {
				klog.V(2).Info(err)
				if res, err = proto.ErrorResponse(msg.Param, err.Error()); err != nil {
					klog.V(2).Info(err)
				}
			}
		case proto.SvcDelete:
			if err = c.handle.KubernetesApi().Delete(msg.Param, msg.Data); err != nil {
				klog.V(2).Info(err)
				if res, err = proto.ErrorResponse(msg.Param, err.Error()); err != nil {
					klog.V(2).Info(err)
				}
			} else {
				if res, err = proto.GetResponse(msg.Param, []byte("delete success")); err != nil {
					klog.V(2).Info(err)
				}
			}
		case proto.SvcGet:
			if res, err = c.handle.KubernetesApi().Get(msg.Param, msg.Data); err != nil {
				klog.V(2).Info(err)
				if res, err = proto.ErrorResponse(msg.Param, err.Error()); err != nil {
					klog.V(2).Info(err)
				}
			}
		case proto.SvcList:
			if res, err = c.handle.KubernetesApi().List(msg.Param); err != nil {
				klog.V(2).Info(err)
				if res, err = proto.ErrorResponse(msg.Param, err.Error()); err != nil {
					klog.V(2).Info(err)
				}
			}
		case proto.SvcWatch:
		case proto.SvcResource:
			if res, err = c.handle.KubernetesApi().Resources(msg.Param); err != nil {
				klog.V(2).Info(err)
				if res, err = proto.ErrorResponse(msg.Param, err.Error()); err != nil {
					klog.V(2).Info(err)
				}
			}
		case proto.SvcHarbor:
			if res, err = c.handle.HarborApi().Core(msg.Param, msg.Data); err != nil {
				klog.V(2).Info(err)
				if res, err = proto.ErrorResponse(msg.Param, err.Error()); err != nil {
					klog.V(2).Info(err)
				}
			}
		}
		if err = c.SendToChannel(res); err != nil {
			return err
		}
	}
}

func (c *wsConn) SendToChannel(msg []byte) (err error) {
	if c.status == connClosed {
		return
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.status == connClosed {
		return
	}
	after := time.After(time.Second * 2)
	for {
		select {
		case c.writeChan <- msg:
			return
		case <-after:
			err = fmt.Errorf("wsSend timeout ws.cid:%d msg:(%v) ws:%v\n", c.clientId, msg, c)
			klog.V(2).Info(err)
			return err
		}
	}
}

func (c *wsConn) WritePump() (err error) {
	defer c.Close()
	for {
		select {
		case <-c.ctx.Done():
			return nil
		case msg, isClose := <-c.writeChan:
			if !isClose {
				return nil
			}
			//klog.Info("send to:", c.clientId, " msg:", string(msg))
			s := time.Now()
			if err := c.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				klog.V(2).Info(err)
				return err
			}
			klog.Info("send to:", c.clientId, " used time:", time.Now().Sub(s))
		}
	}
}

func (c *wsConn) Close() {
	if c.status == connClosed {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.status == connClosed {
		return
	}
	c.status = connClosed
	if err := c.conn.Close(); err != nil {
		klog.V(2).Info(err)
	}
	close(c.writeChan)
	close(c.readChan)
}

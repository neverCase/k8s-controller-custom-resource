package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog"
)

type ConnHub interface {
	NewConn(conn *websocket.Conn)
	GetClientId() int32
	Close()
}

type connHub struct {
	mu           sync.RWMutex
	autoClientId int32
	connections  map[int32]Conn
	ctx          context.Context
}

func (ch *connHub) GetClientId() int32 {
	return atomic.AddInt32(&ch.autoClientId, 1)
}

func (ch *connHub) NewConn(conn *websocket.Conn) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	id := ch.GetClientId()
	ch.connections[id] = NewConn(id, ch.ctx, conn)
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
}

func (ch *connHub) Close() {
}

func NewConnHub(ctx context.Context) ConnHub {
	return &connHub{
		autoClientId: 0,
		connections:  make(map[int32]Conn, 0),
		ctx:          ctx,
	}
}

type Conn interface {
	Ping()
	ReadPump() (err error)
	SendToChannel(data interface{}) (err error)
	WritePump() (err error)
	Close()
}

func NewConn(clientId int32, ctx context.Context, ws *websocket.Conn) Conn {
	c := &conn{
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

func (c *conn) ReadPump() (err error) {
	defer c.Close()
	for {
		var msg Request
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
		switch msg.Service {
		case SvcPing:
			if err = c.SendToChannel(GetResponse(msg.Data)); err != nil {
				return err
			}
		case SvcList:
		case SvcWatch:
		case SvcAdd:
		case SvcUpdate:
		case SvcDelete:
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

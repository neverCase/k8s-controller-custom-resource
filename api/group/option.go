package group

import (
	"fmt"
	"reflect"
	"sync"

	mysqlOperatorV1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
)

type ResourceType string

const (
	ConfigMap   ResourceType = "ConfigMap"
	Deployment  ResourceType = "Deployment"
	Pod         ResourceType = "Pod"
	Secret      ResourceType = "Secret"
	Service     ResourceType = "Service"
	StatefulSet ResourceType = "StatefulSet"

	MysqlOperator ResourceType = "MysqlOperator"
	RedisOperator ResourceType = "RedisOperator"
)

type Options interface {
	Add(opts ...Option)
	Get(rt ResourceType) (Option, error)
	GetReflectType(rt ResourceType) reflect.Type
}

type options struct {
	mu    sync.RWMutex
	hub   map[ResourceType]Option
	kinds map[ResourceType]reflect.Type
}

func (opts *options) Add(opt ...Option) {
	opts.mu.Lock()
	defer opts.mu.Unlock()
	for _, v := range opt {
		if _, ok := opts.kinds[v.Name()]; !ok {
			opts.hub[v.Name()] = v
		}
	}
}

func (opts *options) Get(rt ResourceType) (Option, error) {
	if t, ok := opts.hub[rt]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("err no ResourceType: %s\n", rt)
}

func (opts *options) GetReflectType(rt ResourceType) reflect.Type {
	return opts.kinds[rt]
}

func NewOptions() Options {
	o := &options{
		hub:   make(map[ResourceType]Option, 0),
		kinds: make(map[ResourceType]reflect.Type, 0),
	}
	var b mysqlOperatorV1.MysqlOperator
	a := reflect.New(reflect.TypeOf(b))
	return o
}

type Option interface {
	Name() ResourceType
	Get() interface{}
}

type option struct {
	name ResourceType
	obj  interface{}
}

func (o *option) Name() ResourceType {
	return o.name
}

func (o *option) Get() interface{} {
	return o.obj
}

func NewOption(name ResourceType, obj interface{}) Option {
	o := &option{
		name: name,
		obj:  obj,
	}
	return o
}

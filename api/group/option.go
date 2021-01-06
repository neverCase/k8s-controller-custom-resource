package group

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
)

type Options interface {
	Add(opts ...Option)
	Get(rt ResourceType) (Option, error)
	GetOptionTypeList() []ResourceType
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

func (opts *options) GetOptionTypeList() []ResourceType {
	res := make([]ResourceType, 0)
	for k := range opts.hub {
		res = append(res, k)
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i] < res[j] {
			return true
		}
		return false
	})
	return res
}

func NewOptions() Options {
	o := &options{
		hub:   make(map[ResourceType]Option, 0),
		kinds: make(map[ResourceType]reflect.Type, 0),
	}
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

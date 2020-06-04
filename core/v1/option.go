package v1

import (
	"fmt"
	"k8s.io/klog"
	"reflect"
	"sync"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
)

const (
	// ErrOptionExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Deployment of the same name already existing.
	ErrOptionExists = "ErrOptionExists"

	ErrOptionKindDoesNotExisted = "ErrOptionKindDoesNotExisted"
)

type Options interface {
	Add(opt Option) error
	Get(objType reflect.Type) Option
	GetWithKindName(kindName string) (opt Option, err error)
	HasSyncedFunc() []func() bool
	List() map[reflect.Type]Option
}

type options struct {
	mu    sync.RWMutex
	hub   map[reflect.Type]Option
	kinds map[string]reflect.Type
}

func NewOptions() Options {
	var o Options = &options{
		hub:   make(map[reflect.Type]Option, 0),
		kinds: make(map[string]reflect.Type, 0),
	}
	return o
}

func (o *options) Add(opt Option) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	t := opt.GetReflectType()
	if _, ok := o.hub[t]; ok {
		return fmt.Errorf("%s type:%v", ErrOptionExists, t)
	}
	o.hub[t] = opt
	o.kinds[opt.KindName()] = opt.GetReflectType()
	klog.Info("opt:", opt)
	klog.Info("opt.hub:", o.hub)
	klog.Info("opt.kinds:", o.kinds)
	return nil
}

func (o *options) Get(objType reflect.Type) Option {
	return o.hub[objType]
}

func (o *options) GetWithKindName(kindName string) (opt Option, err error) {
	t, ok := o.kinds[kindName]
	if !ok {
		return opt, fmt.Errorf("%s kind:%v", ErrOptionKindDoesNotExisted, kindName)
	}
	return o.Get(t), nil
}

func (o *options) HasSyncedFunc() []func() bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	res := make([]func() bool, 0)
	for _, v := range o.hub {
		res = append(res, v.HasSyncedFunc())
	}
	return res
}

func (o *options) List() map[reflect.Type]Option {
	return o.hub
}

type Option interface {
	GetReflectType() reflect.Type
	KindName() string
	AgentName() string
	SyncHandleObject(obj interface{}, ks KubernetesResource, recorder record.EventRecorder) error
	HasSyncedFunc() func() bool
	AddEventHandler(handler cache.ResourceEventHandler)
	CompareResourceVersion(old, new interface{}) bool
	Get(nameSpace, ownerRefName string) (obj interface{}, err error)
}

type option struct {
	operatorType               reflect.Type
	kindName                   string
	agentClientSet             interface{}
	agent                      interface{}
	agentName                  string
	hasSynced                  func() bool
	addEvent                   func(handler cache.ResourceEventHandler)
	compareResourceVersionFunc func(old, new interface{}) bool
	getFunc                    func(informer interface{}, nameSpace, ownerRefName string) (obj interface{}, err error)
	syncFunc                   func(obj interface{}, agentClientSet interface{}, ks KubernetesResource, opt record.EventRecorder) error
}

func NewOption(operator interface{},
	agentName, kindName string,
	err error,
	agentClientSet interface{},
	foo interface{},
	hasSynced func() bool,
	addEvent func(handler cache.ResourceEventHandler),
	compareResourceVersionFunc func(old, new interface{}) bool,
	getFunc func(informer interface{}, nameSpace, ownerRefName string) (obj interface{}, err error),
	syncFunc func(obj interface{}, agentClientSet interface{}, ks KubernetesResource, opt record.EventRecorder) error) Option {

	utilruntime.Must(err)

	var opt Option = &option{
		operatorType:               reflect.TypeOf(operator),
		kindName:                   kindName,
		agentClientSet:             agentClientSet,
		agent:                      foo,
		agentName:                  agentName,
		hasSynced:                  hasSynced,
		addEvent:                   addEvent,
		compareResourceVersionFunc: compareResourceVersionFunc,
		getFunc:                    getFunc,
		syncFunc:                   syncFunc,
	}
	return opt
}

func (opt *option) GetReflectType() reflect.Type {
	return opt.operatorType
}

func (opt *option) KindName() string {
	return opt.kindName
}

func (opt *option) AgentName() string {
	return opt.agentName
}

func (opt *option) HasSyncedFunc() func() bool {
	return opt.hasSynced
}

func (opt *option) AddEventHandler(handler cache.ResourceEventHandler) {
	opt.addEvent(handler)
}

func (opt *option) SyncHandleObject(obj interface{}, ks KubernetesResource, recorder record.EventRecorder) error {
	return opt.syncFunc(obj, opt.agentClientSet, ks, recorder)
}

func (opt *option) CompareResourceVersion(old, new interface{}) bool {
	return opt.compareResourceVersionFunc(old, new)
}

func (opt *option) Get(nameSpace, ownerRefName string) (obj interface{}, err error) {
	return opt.getFunc(opt.agent, nameSpace, ownerRefName)
}

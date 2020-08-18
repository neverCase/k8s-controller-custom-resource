package v1

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
)

const (
	// ErrOptionExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Deployment of the same name already existing.
	ErrOptionExists = "ErrOptionExists"

	ErrOptionKindDoesNotExisted = "ErrOptionKindDoesNotExisted"

	ErrOptionWriteWatchChanTimeout = "ErrOptionWriteWatchChanTimeout"
)

type Options interface {
	Add(opt ...Option) error
	Get(objType reflect.Type) Option
	GetWithKindName(kindName string) (opt Option, err error)
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

func (o *options) Add(opt ...Option) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, v := range opt {
		t := v.GetReflectType()
		if _, ok := o.hub[t]; ok {
			return fmt.Errorf("%s type:%v", ErrOptionExists, t)
		}
		o.hub[t] = v
		o.kinds[v.KindName()] = v.GetReflectType()
	}
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

func (o *options) List() map[reflect.Type]Option {
	return o.hub
}

type Option interface {
	GetReflectType() reflect.Type
	KindName() string
	AgentName() string
	Informer() cache.SharedIndexInformer
	SyncHandleObject(obj interface{}, ks KubernetesResource, recorder record.EventRecorder) error
	CompareResourceVersion(old, new interface{}) bool
	Get(nameSpace, ownerRefName string) (obj interface{}, err error)
	SyncObjectStatus(obj interface{}, ks KubernetesResource, recorder record.EventRecorder) error
	WriteWatchChan(e watch.Event, ks KubernetesResource, recorder record.EventRecorder) (err error)
	Watch()
}

type option struct {
	operatorType               reflect.Type
	kindName                   string
	agentClientSet             interface{}
	agent                      interface{}
	agentName                  string
	informer                   cache.SharedIndexInformer
	compareResourceVersionFunc func(old, new interface{}) bool
	getFunc                    func(informer interface{}, nameSpace, ownerRefName string) (obj interface{}, err error)
	syncFunc                   func(obj interface{}, agentClientSet interface{}, ks KubernetesResource, opt record.EventRecorder) error
	syncStatusFunc             func(obj interface{}, agentClientSet interface{}, ks KubernetesResource, recorder record.EventRecorder) error

	watchChan chan OptionWatch
}

func NewOption(operator interface{},
	agentName, kindName string,
	err error,
	agentClientSet interface{},
	foo interface{},
	informer cache.SharedIndexInformer,
	compareResourceVersionFunc func(old, new interface{}) bool,
	getFunc func(informer interface{}, nameSpace, ownerRefName string) (obj interface{}, err error),
	syncFunc func(obj interface{}, agentClientSet interface{}, ks KubernetesResource, opt record.EventRecorder) error,
	syncStatusFunc func(obj interface{}, agentClientSet interface{}, ks KubernetesResource, recorder record.EventRecorder) error) Option {

	utilruntime.Must(err)

	var opt = &option{
		operatorType:               reflect.TypeOf(operator),
		kindName:                   kindName,
		agentClientSet:             agentClientSet,
		agent:                      foo,
		agentName:                  agentName,
		informer:                   informer,
		compareResourceVersionFunc: compareResourceVersionFunc,
		getFunc:                    getFunc,
		syncFunc:                   syncFunc,
		syncStatusFunc:             syncStatusFunc,
		watchChan:                  make(chan OptionWatch, 4096),
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

func (opt *option) Informer() cache.SharedIndexInformer {
	return opt.informer
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

func (opt *option) SyncObjectStatus(obj interface{}, ks KubernetesResource, recorder record.EventRecorder) error {
	return opt.syncStatusFunc(obj, opt.agentClientSet, ks, recorder)
}

type OptionWatch struct {
	Resource KubernetesResource
	Recorder record.EventRecorder
	Event    watch.Event
}

func (opt *option) WriteWatchChan(e watch.Event, ks KubernetesResource, recorder record.EventRecorder) (err error) {
	after := time.After(time.Millisecond * 500)
	for {
		select {
		case <-after:
			return fmt.Errorf("%s kind:%v", ErrOptionWriteWatchChanTimeout, opt.kindName)
		case opt.watchChan <- OptionWatch{Event: e, Resource: ks, Recorder: recorder}:
		}
	}
}

func (opt *option) Watch() {
	for {
		select {
		case ow, isClosed := <-opt.watchChan:
			if !isClosed {
				return
			}
			if err := opt.syncStatusFunc(ow.Event, opt.agentClientSet, ow.Resource, ow.Recorder); err != nil {
				klog.V(2).Info(err)
			}
		}
	}
}

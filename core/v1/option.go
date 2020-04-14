package v1

import (
	"reflect"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
)

type KubernetesOption interface {
}

type kubernetesOption struct {
	recorder                   record.EventRecorder
	kindName                   string
	agentClientSet             interface{}
	agent                      interface{}
	agentName                  string
	agentType                  reflect.Type
	hasSynced                  func() bool
	addEvent                   func(handler cache.ResourceEventHandler)
	compareResourceVersionFunc func(old, new interface{}) bool
	getFunc                    func(informer interface{}, nameSpace, ownerRefName string) (obj interface{}, err error)
	syncFunc                   func(obj interface{}, agentClientSet interface{}, ks KubernetesResource, ko record.EventRecorder) error
}

func NewKubernetesOption(agentName, kindName string,
	agentClientSet interface{},
	foo interface{},
	hasSynced func() bool,
	addEvent func(handler cache.ResourceEventHandler),
	compareResourceVersionFunc func(old, new interface{}) bool,
	getFunc func(informer interface{}, nameSpace, ownerRefName string) (obj interface{}, err error),
	syncFunc func(obj interface{}, agentClientSet interface{}, ks KubernetesResource, ko record.EventRecorder) error) KubernetesOption {

	var opt KubernetesOption = &kubernetesOption{
		kindName:                   kindName,
		agentClientSet:             agentClientSet,
		agent:                      foo,
		agentName:                  agentName,
		agentType:                  reflect.TypeOf(foo),
		hasSynced:                  hasSynced,
		addEvent:                   addEvent,
		compareResourceVersionFunc: compareResourceVersionFunc,
		getFunc:                    getFunc,
		syncFunc:                   syncFunc,
	}
	return opt
}

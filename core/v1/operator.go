package v1

import (
	"reflect"
	"time"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type KubernetesOperator interface {
	GetClientSet() kubernetes.Interface
	GetInformerFactory() kubeinformers.SharedInformerFactory
	GetKindName() string
	GetAgentName() string
	SyncHandleObject(obj interface{}) error
	HandleObject(obj interface{})
	HasSyncedFunc() func() bool
	AddEventHandler(handler cache.ResourceEventHandler)
	CompareResourceVersion(old, new interface{}) bool
	Get(nameSpace, ownerRefName string) (obj interface{}, err error)
}

func NewKubernetesOperator(kubeClientset kubernetes.Interface,
	agentName, kindName string,
	foo interface{},
	hasSynced func() bool,
	addEvent func(handler cache.ResourceEventHandler),
	compareResourceVersionFunc func(old, new interface{}) bool,
	getFunc func(informer interface{}, nameSpace, ownerRefName string) (obj interface{}, err error),
	syncFunc func(obj interface{}) error) KubernetesOperator {
	var ko KubernetesOperator = &kubernetesOperator{
		kubeClientSet:              kubeClientset,
		kubeInformerFactory:        kubeinformers.NewSharedInformerFactory(kubeClientset, time.Second*30),
		kindName:                   kindName,
		agent:                      foo,
		agentName:                  agentName,
		agentType:                  reflect.TypeOf(foo),
		hasSynced:                  hasSynced,
		addEvent:                   addEvent,
		compareResourceVersionFunc: compareResourceVersionFunc,
		getFunc:                    getFunc,
		syncFunc:                   syncFunc,
	}
	return ko
}

type kubernetesOperator struct {
	kubeClientSet              kubernetes.Interface
	kubeInformerFactory        kubeinformers.SharedInformerFactory
	kindName                   string
	agent                      interface{}
	agentName                  string
	agentType                  reflect.Type
	hasSynced                  func() bool
	addEvent                   func(handler cache.ResourceEventHandler)
	compareResourceVersionFunc func(old, new interface{}) bool
	getFunc                    func(informer interface{}, nameSpace, ownerRefName string) (obj interface{}, err error)
	syncFunc                   func(obj interface{}) error
}

func (ko *kubernetesOperator) GetClientSet() kubernetes.Interface {
	return ko.kubeClientSet
}

func (ko *kubernetesOperator) GetInformerFactory() kubeinformers.SharedInformerFactory {
	return ko.kubeInformerFactory
}

func (ko *kubernetesOperator) GetKindName() string {
	return ko.kindName
}

func (ko *kubernetesOperator) GetAgentName() string {
	return ko.agentName
}

func (ko *kubernetesOperator) HasSyncedFunc() func() bool {
	return ko.hasSynced
}

func (ko *kubernetesOperator) AddEventHandler(handler cache.ResourceEventHandler) {
	ko.addEvent(handler)
}

func (ko *kubernetesOperator) HandleObject(obj interface{}) {

}

func (ko *kubernetesOperator) SyncHandleObject(obj interface{}) error {
	return ko.syncFunc(obj)
}

func (ko *kubernetesOperator) CompareResourceVersion(old, new interface{}) bool {
	return ko.compareResourceVersionFunc(old, new)
}

func (ko *kubernetesOperator) Get(nameSpace, ownerRefName string) (obj interface{}, err error) {
	return ko.getFunc(ko.agent, nameSpace, ownerRefName)
}

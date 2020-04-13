package v1

import (
	"reflect"
	"time"

	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
)

type KubernetesOperator interface {
	GetClientSet() kubernetes.Interface
	GetRecorder() record.EventRecorder
	GetInformerFactory() kubeinformers.SharedInformerFactory
	GetResource() KubernetesResource
	GetKindName() string
	GetAgentName() string
	SyncHandleObject(obj interface{}) error
	HasSyncedFunc() func() bool
	AddEventHandler(handler cache.ResourceEventHandler)
	CompareResourceVersion(old, new interface{}) bool
	Get(nameSpace, ownerRefName string) (obj interface{}, err error)
}

func NewKubernetesOperator(kubeClientset kubernetes.Interface,
	stopCh <-chan struct{},
	err error,
	agentName, kindName string,
	agentClientSet interface{},
	foo interface{},
	hasSynced func() bool,
	addEvent func(handler cache.ResourceEventHandler),
	compareResourceVersionFunc func(old, new interface{}) bool,
	getFunc func(informer interface{}, nameSpace, ownerRefName string) (obj interface{}, err error),
	syncFunc func(obj interface{}, agentClientSet interface{}, ks KubernetesResource, ko record.EventRecorder) error) KubernetesOperator {

	utilruntime.Must(err)
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: agentName})

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClientset, time.Second*30)

	var ko KubernetesOperator = &kubernetesOperator{
		kubeClientSet:              kubeClientset,
		kubeInformerFactory:        kubeInformerFactory,
		recorder:                   recorder,
		kubernetesResource:         NewKubernetesResource(kubeClientset, kubeInformerFactory, recorder),
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

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine
	kubeInformerFactory.Start(stopCh)
	return ko
}

type kubernetesOperator struct {
	kubernetesResource         KubernetesResource
	kubeClientSet              kubernetes.Interface
	kubeInformerFactory        kubeinformers.SharedInformerFactory
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

func (ko *kubernetesOperator) GetClientSet() kubernetes.Interface {
	return ko.kubeClientSet
}

func (ko *kubernetesOperator) GetRecorder() record.EventRecorder {
	return ko.recorder
}

func (ko *kubernetesOperator) GetInformerFactory() kubeinformers.SharedInformerFactory {
	return ko.kubeInformerFactory
}

func (ko *kubernetesOperator) GetResource() KubernetesResource {
	return ko.kubernetesResource
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

func (ko *kubernetesOperator) SyncHandleObject(obj interface{}) error {
	return ko.syncFunc(obj, ko.agentClientSet, ko.kubernetesResource, ko.recorder)
}

func (ko *kubernetesOperator) CompareResourceVersion(old, new interface{}) bool {
	return ko.compareResourceVersionFunc(old, new)
}

func (ko *kubernetesOperator) Get(nameSpace, ownerRefName string) (obj interface{}, err error) {
	return ko.getFunc(ko.agent, nameSpace, ownerRefName)
}

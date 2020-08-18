package v1

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
)

type KubernetesOperator interface {
	Recorder() record.EventRecorder
	InformerFactory() kubeinformers.SharedInformerFactory
	Resource() KubernetesResource
	AgentName() string
	Options() Options
	Watch()
}

func NewKubernetesOperator(kubeClientset kubernetes.Interface,
	stopCh <-chan struct{},
	agentName string,
	opts Options) KubernetesOperator {

	//utilruntime.Must(err)
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: agentName})

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClientset, time.Second*30)

	var ko = &kubernetesOperator{
		kubeClientSet:       kubeClientset,
		kubeInformerFactory: kubeInformerFactory,
		recorder:            recorder,
		kubernetesResource:  NewKubernetesResource(kubeClientset, kubeInformerFactory),
		agentName:           agentName,
		options:             opts,
	}

	ko.Watch()

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine
	kubeInformerFactory.Start(stopCh)
	return ko
}

type kubernetesOperator struct {
	kubernetesResource KubernetesResource
	// kubeclientset is a standard kubernetes clientset
	kubeClientSet       kubernetes.Interface
	kubeInformerFactory kubeinformers.SharedInformerFactory
	recorder            record.EventRecorder
	agentName           string
	options             Options
}

func (ko *kubernetesOperator) Recorder() record.EventRecorder {
	return ko.recorder
}

func (ko *kubernetesOperator) InformerFactory() kubeinformers.SharedInformerFactory {
	return ko.kubeInformerFactory
}

func (ko *kubernetesOperator) Resource() KubernetesResource {
	return ko.kubernetesResource
}

func (ko *kubernetesOperator) AgentName() string {
	return ko.agentName
}

func (ko *kubernetesOperator) Options() Options {
	return ko.options
}

const LabelsFilterNameTemplate = "app=%s"

func (ko *kubernetesOperator) Watch() {
	for _, opt := range ko.Options().List() {
		go func(opt Option) {
			res, err := ko.kubernetesResource.StatefulSet().Watch("", fmt.Sprintf(LabelsFilterNameTemplate, opt.KindName()))
			if err != nil {
				klog.Fatal(err)
			}
			for {
				select {
				case e, isClosed := <-res.ResultChan():
					if !isClosed {
						return
					}
					if err := opt.WriteWatchChan(e, ko.kubernetesResource, ko.recorder); err != nil {
						klog.V(2).Info(err)
					}
				}
			}
		}(opt)
	}
}

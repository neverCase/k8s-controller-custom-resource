package v1

import (
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
)

type KubernetesResource interface {
	Deployment() KubernetesDeployment
	Service() KubernetesService
	StatefulSet() KubernetesStatefulSet
	ConfigMap() KubernetesConfigMap
}

type kubernetesResource struct {
	kubeClientSet       kubernetes.Interface
	kubeInformerFactory kubeinformers.SharedInformerFactory
	recorder            record.EventRecorder

	deployment  KubernetesDeployment
	service     KubernetesService
	statefulSet KubernetesStatefulSet
	configMap   KubernetesConfigMap
}

func NewKubernetesResource(kubeClientSet kubernetes.Interface, kubeInformerFactory kubeinformers.SharedInformerFactory, recorder record.EventRecorder) KubernetesResource {
	var kd KubernetesResource = &kubernetesResource{
		kubeClientSet:       kubeClientSet,
		kubeInformerFactory: kubeInformerFactory,
		recorder:            recorder,
		deployment:          NewKubernetesDeployment(kubeClientSet, kubeInformerFactory, recorder),
		service:             NewKubernetesService(kubeClientSet, kubeInformerFactory, recorder),
		statefulSet:         NewKubernetesStatefulSet(kubeClientSet, kubeInformerFactory, recorder),
		configMap:           NewKubernetesConfigMap(kubeClientSet, kubeInformerFactory, recorder),
	}
	return kd
}

func (kr *kubernetesResource) Deployment() KubernetesDeployment {
	return kr.deployment
}

func (kr *kubernetesResource) Service() KubernetesService {
	return kr.service
}

func (kr *kubernetesResource) StatefulSet() KubernetesStatefulSet {
	return kr.statefulSet
}

func (kr *kubernetesResource) ConfigMap() KubernetesConfigMap {
	return kr.configMap
}

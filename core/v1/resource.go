package v1

import (
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

type KubernetesResource interface {
	ClientSet() kubernetes.Interface
	Deployment() KubernetesDeployment
	Service() KubernetesService
	StatefulSet() KubernetesStatefulSet
	ConfigMap() KubernetesConfigMap
}

type kubernetesResource struct {
	kubeClientSet       kubernetes.Interface
	kubeInformerFactory kubeinformers.SharedInformerFactory

	deployment  KubernetesDeployment
	service     KubernetesService
	statefulSet KubernetesStatefulSet
	configMap   KubernetesConfigMap
}

func NewKubernetesResource(kubeClientSet kubernetes.Interface, kubeInformerFactory kubeinformers.SharedInformerFactory) KubernetesResource {
	var kr KubernetesResource = &kubernetesResource{
		kubeClientSet:       kubeClientSet,
		kubeInformerFactory: kubeInformerFactory,
		deployment:          NewKubernetesDeployment(kubeClientSet, kubeInformerFactory),
		service:             NewKubernetesService(kubeClientSet, kubeInformerFactory),
		statefulSet:         NewKubernetesStatefulSet(kubeClientSet, kubeInformerFactory),
		configMap:           NewKubernetesConfigMap(kubeClientSet, kubeInformerFactory),
	}
	return kr
}

func (kr *kubernetesResource) ClientSet() kubernetes.Interface {
	return kr.kubeClientSet
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

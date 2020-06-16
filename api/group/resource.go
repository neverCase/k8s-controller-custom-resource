package group

import (
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
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
	HelixOperator ResourceType = "HelixOperator"
)

type Resource interface {
}

type resource struct {
	kubeClientSet       kubernetes.Interface
	kubeInformerFactory kubeinformers.SharedInformerFactory
}

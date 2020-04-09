package v1

import "k8s.io/client-go/tools/cache"

type KubernetesInformer interface {
	Informer() cache.SharedIndexInformer
}

type kubernetsInformer struct {
}

func (ki *kubernetsInformer) Informer() {

}

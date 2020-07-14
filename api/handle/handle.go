package handle

import "github.com/nevercase/k8s-controller-custom-resource/api/group"

type Handle interface {
	KubernetesApiGetter
	HarborApiGetter
}

func NewHandle(g group.Group) Handle {
	return &handle{
		group:  g,
		k8s:    NewKubernetesApiHandle(g),
		harbor: NewHandle(g),
	}
}

type handle struct {
	group  group.Group
	k8s    KubernetesApiInterface
	harbor HarborApiInterface
}

func (h *handle) KubernetesApi() KubernetesApiInterface {
	return h.k8s
}

func (h *handle) HarborApi() HarborApiInterface {
	return h.harbor
}

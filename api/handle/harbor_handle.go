package handle

import "github.com/nevercase/k8s-controller-custom-resource/api/group"

type HarborApiGetter interface {
	HarborApi() HarborApiInterface
}

type HarborApiInterface interface {
}

type harborApi struct {
	group group.Group
}

func (ha *harborApi) Projects() {

}

package group

import harbor "github.com/nevercase/harbor-api"

type Group interface {
	ResourceGetter
	harbor.HubGetter
}

func NewGroup(masterUrl, kubeConfigPath string, dockerHub []harbor.Config) Group {
	var g = &group{
		resource: NewResource(masterUrl, kubeConfigPath),
		harbor:   harbor.NewHub(dockerHub),
	}
	return g
}

type group struct {
	resource ResourceInterface
	harbor   harbor.HubInterface
}

func (g *group) Resource() ResourceInterface {
	return g.resource
}

func (g *group) HarborHub() harbor.HubInterface {
	return g.harbor
}

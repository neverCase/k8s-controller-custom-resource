package group

import harbor "github.com/nevercase/harbor-api"

type Group interface {
	ResourceGetter
	harbor.HarborGetter
}

func NewGroup(masterUrl, kubeconfigPath, dockerUrl, dockerAdmin, dockerPassword string) Group {
	var g = &group{
		resource: NewResource(masterUrl, kubeconfigPath),
		harbor:   harbor.NewHarbor(dockerUrl, dockerAdmin, dockerPassword),
	}
	return g
}

type group struct {
	resource ResourceInterface
	harbor   harbor.HarborInterface
}

func (g *group) Resource() ResourceInterface {
	return g.resource
}

func (g *group) Harbor() harbor.HarborInterface {
	return g.harbor
}

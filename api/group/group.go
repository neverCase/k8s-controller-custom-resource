package group

type Group interface {
	ResourceGetter
	HarborGetter
}

func NewGroup(masterUrl, kubeconfigPath, dockerUrl, dockerAdmin, dockerPassword string) Group {
	var g = &group{
		resource: NewResource(masterUrl, kubeconfigPath),
		harbor:   NewHarbor(dockerUrl, dockerAdmin, dockerPassword),
	}
	return g
}

type group struct {
	resource ResourceInterface
	harbor   HarborInterface
}

func (g *group) Resource() ResourceInterface {
	return g.resource
}

func (g *group) Harbor() HarborInterface {
	return g.harbor
}

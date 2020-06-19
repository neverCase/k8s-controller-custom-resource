package group

type Group interface {
	ResourceGetter
}

func NewGroup(masterUrl, kubeconfigPath string) Group {
	var g = &group{
		resource: NewResource(masterUrl, kubeconfigPath),
	}
	return g
}

type group struct {
	resource ResourceInterface
}

func (g *group) Resource() ResourceInterface {
	return g.resource
}

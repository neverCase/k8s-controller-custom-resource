package group

import (
	"context"

	harbor "github.com/nevercase/harbor-api"
	"k8s.io/apimachinery/pkg/watch"
)

type Group interface {
	ResourceGetter
	harbor.HubGetter

	WatchEvents() chan watch.Event
}

func NewGroup(ctx context.Context, masterUrl, kubeConfigPath string, dockerHub []harbor.Config) Group {
	events := make(chan watch.Event, 1024)
	var g = &group{
		resource: NewResource(ctx, masterUrl, kubeConfigPath, events),
		harbor:   harbor.NewHub(dockerHub),
		events:   events,
	}
	return g
}

type group struct {
	resource ResourceInterface
	harbor   harbor.HubInterface

	events chan watch.Event
}

func (g *group) Resource() ResourceInterface {
	return g.resource
}

func (g *group) HarborHub() harbor.HubInterface {
	return g.harbor
}

func (g *group) WatchEvents() chan watch.Event {
	return g.events
}

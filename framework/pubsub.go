package framework

import (
	"github.com/kgrunwald/goweb/di"
	"github.com/kgrunwald/goweb/ilog"
	"github.com/kgrunwald/goweb/pubsub"
)

type pubSubDef struct {
	Handlers []string
}

// InitializePubSub loads in the configuration of any PubSub handlers and registers them with the PubSub Bus
func InitializePubSub(bus pubsub.Bus, container di.Container, logger ilog.Logger) {
	bindings := loadPubSubYaml()
	for _, binding := range bindings {
		logger.WithFields("binding", binding).
			Debug("Adding PubSub Handler")
		m := container.GetMethod(binding.Service(), binding.Method)
		bus.Subscribe(m.Interface())
	}
}

func loadPubSubYaml() []Binding {
	pubSub := pubSubDef{}
	LoadYaml("pubsub.yaml", &pubSub)

	bindings := []Binding{}
	for _, handler := range pubSub.Handlers {
		bindings = append(bindings, NewBinding(handler))
	}
	return bindings
}

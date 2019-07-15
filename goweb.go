package goweb

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/kgrunwald/goweb/di"
	"github.com/kgrunwald/goweb/ilog"
	"github.com/kgrunwald/goweb/pubsub"
	"github.com/kgrunwald/goweb/router"
)

func init() {
	c := di.GetContainer()
	c.Invoke(func(log ilog.Logger) {
		if err := godotenv.Load(); err != nil {
			log.Debug("No .env file found")
		}
	})
}

// Start initializes all of the dependencies of the framework and starts listening for incoming HTTP requests
func Start() {
	c := di.GetContainer()
	c.Invoke(func(r router.Router, log ilog.Logger) {
		log.Info(fmt.Sprintf("Listening on port %d", 80))
		r.Start(80)
	})
}

func Route(path string, method interface{}) router.Route {
	var route router.Route
	di.GetContainer().Invoke(func(r router.Router, log ilog.Logger) {
		route = r.NewRoute()
		route.Path(path)

		m := reflect.ValueOf(method)
		vars := []string{}
		re := regexp.MustCompile(`\{([^{}]+)\}`)
		matches := re.FindAllStringSubmatch(route.GetPath(), -1)
		for _, match := range matches {
			vars = append(vars, match[1])
		}

		binding := router.RouteBinding{
			Route: route,
			Vars:  vars,
		}
		handler := &router.RouteHandler{Method: m, Binding: binding, Router: r, Log: log}
		route.Handler(handler.Handle)

		log.WithFields(
			"Route", path).
			Debug("Adding route")
	})

	return route
}

func Subscribe(method interface{}) {
	di.GetContainer().Invoke(func(bus pubsub.Bus, logger ilog.Logger) {
		logger.Debug("Adding PubSub handler")
		bus.Subscribe(method)
	})
}

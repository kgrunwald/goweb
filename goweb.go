package goweb

import (
	"fmt"
	"os"
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
	c.Register(getServerInfo)
	c.Invoke(func(log ilog.Logger) {
		if err := godotenv.Load(); err != nil {
			log.Debug("No .env file found")
		}
	})
}

type serverInfo struct {
	port int
}

func getServerInfo(log ilog.Logger) *serverInfo {
	port, err := strconf.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Failed to get $PORT environment variable")
	}

	return &serverInfo{port}
}

// Start initializes all of the dependencies of the framework and starts listening for incoming HTTP requests
func Start() {
	c := di.GetContainer()
	c.Invoke(func(r router.Router, log ilog.Logger, info *serverInfo) {
		log.Info(fmt.Sprintf("Listening on port %d", info.port))
		r.Start(info.port)
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

func ServeSPA(pathPrefix, staticPath string) {
	di.GetContainer().Invoke(func(r router.Router, log ilog.Logger) {
		r.ServeSPA(pathPrefix, staticPath)
	})
}

func Subscribe(method interface{}) {
	di.GetContainer().Invoke(func(bus pubsub.Bus, logger ilog.Logger) {
		logger.Debug("Adding PubSub handler")
		bus.Subscribe(method)
	})
}

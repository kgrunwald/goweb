package goweb

import (
	"fmt"
	"os"
	"strconv"

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

type ServerInfo struct {
	Port int
}

func (s *ServerInfo) String() string {
	return fmt.Sprintf("localhost:%d", s.Port)
}

func getServerInfo(log ilog.Logger) *ServerInfo {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Failed to get $PORT environment variable")
	}

	return &ServerInfo{port}
}

// Start initializes all of the dependencies of the framework and starts listening for incoming HTTP requests
func Start() {
	c := di.GetContainer()
	c.Invoke(func(r router.Router, log ilog.Logger, info *ServerInfo) {
		log.Info(fmt.Sprintf("Listening on port %d", info.Port))
		r.Start(info.Port)
	})
}

func Route(path string, method interface{}) router.Route {
	var route router.Route
	di.GetContainer().Invoke(func(r router.Router, log ilog.Logger) {
		route = r.Route(path, method)
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

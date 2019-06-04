package framework

import (
	"log"

	"github.com/kgrunwald/goweb/di"
	"github.com/kgrunwald/goweb/router"

	"github.com/joho/godotenv"
)

// Start initializes all of the dependencies of the framework and starts listening for incoming HTTP requests
func Start() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	c := di.GetContainer()
	c.Invoke(InitializeRouter)
	c.Invoke(InitializePubSub)
	c.Invoke(func(r router.Router) { r.Start(80) })
}

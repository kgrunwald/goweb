package framework

import (
	"log"

	"github.com/kgrunwald/goweb/di"
	"github.com/kgrunwald/goweb/router"

	"github.com/joho/godotenv"
)

func Start() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	c := di.GetContainer()
	c.Invoke(Initialize)
	c.Invoke(func(r *router.Router) { r.Start() })
}

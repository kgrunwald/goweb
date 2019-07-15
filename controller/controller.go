package controller

import (
	"github.com/kgrunwald/goweb"
	"github.com/kgrunwald/goweb/di"
)

func Register() {
	c := di.GetContainer()
	c.Register(NewT)
	c.Invoke(InitRoutes)
}

func InitRoutes(t *T) {
	goweb.Route("/add/{a}/{b}", t.Add).
		Methods("GET").
		Name("test_route")

	goweb.Route("/add", t.AddPost).
		Name("test_route_post").
		Methods("POST")

	goweb.Route("/", t.GetVersion).
		Methods("POST").
		Headers("SOAPAction", "getVersionInformation")

	goweb.Subscribe(t.MessageHandler)
}

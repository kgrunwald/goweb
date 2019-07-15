package main

import (
	"github.com/kgrunwald/goweb"
	"github.com/kgrunwald/goweb/controller"
)

func main() {
	controller.Register()
	goweb.Start()
}

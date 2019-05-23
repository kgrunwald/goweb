package main

import (
	"github.com/kgrunwald/goweb/controller"
	"github.com/kgrunwald/goweb/framework"
)

func main() {
	controller.Register()
	framework.Start()
}

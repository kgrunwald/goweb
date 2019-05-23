package controller

import "github.com/kgrunwald/goweb/di"

func Register() {
	c := di.GetContainer()
	c.RegisterGroup(NewT, di.GroupController)
}

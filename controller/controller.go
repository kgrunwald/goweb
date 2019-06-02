package controller

import (
	"github.com/kgrunwald/goweb/di"
)

func Register() {
	c := di.GetContainer()
	c.Register(NewT)
}

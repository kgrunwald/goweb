package di

import (
	"fmt"
	"log"

	"go.uber.org/dig"
)

type Group string

const (
	GroupController Group = "controller"
)

type Container interface {
	Register(interface{})
	RegisterName(interface{}, string)
	RegisterGroup(interface{}, Group)
	Invoke(interface{})
	GetControllers() []interface{}
	Print()
}

type Controllers struct {
	dig.In
	Controllers []interface{} `group:"controller"`
}

type container struct {
	digContainer *dig.Container
	controllers  []interface{}
}

var c *container

func init() {
	c = &container{
		digContainer: dig.New(),
	}

	c.Register(c.ResolveContainer)
}

func GetContainer() Container {
	return c
}

func (c *container) ResolveContainer(ctrls Controllers) Container {
	c.controllers = ctrls.Controllers
	return c
}

func (c *container) GetControllers() []interface{} {
	return c.controllers
}

func (c *container) Register(constructor interface{}) {
	err := c.digContainer.Provide(constructor)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *container) RegisterGroup(constructor interface{}, group Group) {
	err := c.digContainer.Provide(constructor, dig.Group(string(group)))
	if err != nil {
		log.Fatal(err)
	}
}

func (c *container) RegisterName(constructor interface{}, name string) {
	err := c.digContainer.Provide(constructor, dig.Name(name))
	if err != nil {
		log.Fatal(err)
	}
}

func (c *container) Invoke(function interface{}) {
	err := c.digContainer.Invoke(function)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *container) Print() {
	fmt.Println(c.digContainer.String())
}

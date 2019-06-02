package di

import (
	"fmt"
	"reflect"
)

type Container interface {
	Register(interface{})
	Get(name string) interface{}
	GetMethod(service, method string) reflect.Value
	Invoke(interface{})
	Print()
}

type ServiceContainer struct {
	Services     map[string]interface{}
	Constructors map[string]Constructor
}

type Constructor func(*ServiceContainer) interface{}

var svc *ServiceContainer

func init() {
	svc = &ServiceContainer{
		Services:     make(map[string]interface{}),
		Constructors: make(map[string]Constructor),
	}

	// add the container as a service in the container
	svc.Register(GetContainer)
}

func GetContainer() Container {
	return svc
}

func (c *ServiceContainer) Register(ctor interface{}) {
	ctorType := reflect.TypeOf(ctor)
	returnType := ctorType.Out(0)
	typeName := c.GetTypeName(returnType)
	c.Constructors[typeName] = func(c *ServiceContainer) interface{} {
		out := c.Call(ctor)
		return out[0].Interface()
	}
}

func (c *ServiceContainer) Invoke(f interface{}) {
	c.Call(f)
}

func (c *ServiceContainer) Call(f interface{}) []reflect.Value {
	fType := reflect.TypeOf(f)
	args := []reflect.Value{}
	for i := 0; i < fType.NumIn(); i++ {
		argType := fType.In(i)
		argName := c.GetTypeName(argType)
		argValue := c.GetValue(argName).Convert(argType)
		args = append(args, argValue)
	}

	return reflect.ValueOf(f).Call(args)
}

func (c *ServiceContainer) GetTypeName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.String()
}

func (c *ServiceContainer) Get(name string) interface{} {
	if _, ok := c.Services[name]; !ok {
		if ctor, ok := c.Constructors[name]; ok {
			c.Services[name] = ctor(c)
		} else {
			panic("Attempted to get service " + name + " from container, but it does not exist.")
		}
	}

	return c.Services[name]
}

func (c *ServiceContainer) GetMethod(service, method string) reflect.Value {
	svc := c.Get(service)
	t := reflect.ValueOf(svc)
	return t.MethodByName(method)
}

func (c *ServiceContainer) GetValue(name string) reflect.Value {
	return reflect.ValueOf(c.Get(name))
}

func (c *ServiceContainer) Print() {
	fmt.Println(c)
}

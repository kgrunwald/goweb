package di

import (
	"fmt"
	"reflect"
)

// The Container interface defines methods for the container
type Container interface {
	Register(interface{})
	Get(name string) interface{}
	GetMethod(service, method string) reflect.Value
	Invoke(interface{})
	Print()
}

// The ServiceContainer is a naive implentation of a container that stores all services as singletons
type ServiceContainer struct {
	Services     map[string]interface{}
	Constructors map[string]Constructor
}

// A Constructor is a function that takes in the container as an argument, retrieves all necessary services, and returns
// a new instance of an `interface{}`
type Constructor func(*ServiceContainer) interface{}

var svc *ServiceContainer

func init() {
	svc = &ServiceContainer{
		Services:     make(map[string]interface{}),
		Constructors: make(map[string]Constructor),
	}

	// Add the container as a service in the container. Meta.
	svc.Register(GetContainer)
}

// GetContainer returns an implementation of a `Container`
func GetContainer() Container {
	return svc
}

// Register adds a service to the container. The provided constructor must only take in other services as arguments,
// and it must return 1 value to be added to the container.
func (c *ServiceContainer) Register(ctor interface{}) {
	ctorType := reflect.TypeOf(ctor)
	returnType := ctorType.Out(0)
	typeName := c.GetTypeName(returnType)
	c.Constructors[typeName] = func(c *ServiceContainer) interface{} {
		out := c.Call(ctor)
		return out[0].Interface()
	}
}

// Invoke will call the provided method with the values resolved from the container. All arguments must be services that
// are defined in the container.
func (c *ServiceContainer) Invoke(f interface{}) {
	c.Call(f)
}

// Call is the same as `Invoke`, but the values returned from the function are propagated back to the caller as `[]reflect.Value`.
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

// GetTypeName generates the name of the service entry in the container from the Type.
func (c *ServiceContainer) GetTypeName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.String()
}

// Get a service from the container by name. If the service is not in the container, the method will panic.
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

// GetMethod looks up the service in the container by name and then looks up the method by name on the returned instance.
func (c *ServiceContainer) GetMethod(service, method string) reflect.Value {
	svc := c.Get(service)
	t := reflect.ValueOf(svc)
	return t.MethodByName(method)
}

// GetValue returns an entry from the container as a reflect.Value
func (c *ServiceContainer) GetValue(name string) reflect.Value {
	return reflect.ValueOf(c.Get(name))
}

// Print prints out the contents of the container
func (c *ServiceContainer) Print() {
	fmt.Println(c)
}

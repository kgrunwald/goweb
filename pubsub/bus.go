package pubsub

import (
	"reflect"

	"github.com/kgrunwald/goweb/di"
	"github.com/kgrunwald/goweb/ilog"
)

// The Bus interface provides a mechanism to subscribe to asynchronous events and to dispatch events to all subscribers.
// Subscribers must pass in a function that takes in a single argument. When the bus dispatches a message whose type or
// interface matches the subscriber, the subscribers function will be invoked. Each invocation will be called in a separate
// goroutine.
type Bus interface {
	Subscribe(interface{})
	Dispatch(data interface{})
}

type queue chan interface{}

func init() {
	container := di.GetContainer()
	container.Register(NewBus)
}

// NewBus returns a concrete implementation of the `Bus` interface
func NewBus(logger ilog.Logger) Bus {
	return newEventBus(logger)
}

type eventBus struct {
	Subscriptions []interface{}
	Queue         queue
	Logger        ilog.Logger
}

func newEventBus(logger ilog.Logger) *eventBus {
	bus := &eventBus{
		Subscriptions: []interface{}{},
		Queue:         make(queue),
		Logger:        logger,
	}

	bus.Run()
	return bus
}

func (e *eventBus) Subscribe(handler interface{}) {
	e.Subscriptions = append(e.Subscriptions, handler)
}

func (e *eventBus) Dispatch(message interface{}) {
	e.Queue <- message
}

func (e *eventBus) Run() {
	go func() {
		for {
			select {
			case message := <-e.Queue:
				go e.Publish(message)
			}
		}
	}()
}

func (e *eventBus) Publish(message interface{}) {
	msgType := reflect.TypeOf(message)
	for _, handler := range e.Subscriptions {
		t := reflect.TypeOf(handler).In(0)
		if t == msgType {
			e.Invoke(handler, message)
		} else if t.Kind() == reflect.Interface && msgType.Implements(t) {
			e.Invoke(handler, message)
		}
	}
}

func (e *eventBus) Invoke(handler, message interface{}) {
	handlerVal := reflect.ValueOf(handler)
	msgVal := reflect.ValueOf(message)
	go handlerVal.Call([]reflect.Value{msgVal})
}

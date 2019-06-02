package pubsub

import (
	"reflect"

	"github.com/kgrunwald/goweb/di"
	"github.com/kgrunwald/goweb/ilog"
)

type Bus interface {
	Subscribe(interface{})
	Dispatch(data interface{})
}

type Queue chan interface{}

func init() {
	container := di.GetContainer()
	container.Register(NewBus)
}

func NewBus(logger ilog.Logger) Bus {
	bus := &EventBus{
		Subscriptions: []interface{}{},
		Queue:         make(Queue),
		Logger:        logger,
	}

	bus.Run()
	return bus
}

type EventBus struct {
	Subscriptions []interface{}
	Queue         Queue
	Logger        ilog.Logger
}

func (e *EventBus) Subscribe(handler interface{}) {
	e.Subscriptions = append(e.Subscriptions, handler)
}

func (e *EventBus) Dispatch(message interface{}) {
	e.Queue <- message
}

func (e *EventBus) Run() {
	go func() {
		for {
			select {
			case message := <-e.Queue:
				go e.Publish(message)
			}
		}
	}()
}

func (e *EventBus) Publish(message interface{}) {
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

func (e *EventBus) Invoke(handler, message interface{}) {
	handlerVal := reflect.ValueOf(handler)
	msgVal := reflect.ValueOf(message)
	go handlerVal.Call([]reflect.Value{msgVal})
}

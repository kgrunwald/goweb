package pubsub

import (
	"time"
)

type Bus interface {
	Subscribe(topic string, handler Handler)
	Dispatch(data Message)
}

type Message struct {
	Topic     string
	Timestamp time.Time
	Payload   interface{}
}

type Queue chan Message

type Handler func(Message)

func NewBus() Bus {
	bus := &EventBus{
		Subscriptions: map[string][]Handler{},
		Queue: make(Queue),
	}

	bus.Run()
	return bus
}

type EventBus struct {
	Subscriptions map[string][]Handler
	Queue         Queue
}

func (e *EventBus) Subscribe(topic string, handler Handler) {
	if _, ok := e.Subscriptions[topic]; !ok {
		e.Subscriptions[topic] = []Handler{}
	}

	e.Subscriptions[topic] = append(e.Subscriptions[topic], handler)
}

func (e *EventBus) Dispatch(message Message) {
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

func (e *EventBus) Publish(message Message) {
	for _, handler := range e.Subscriptions[message.Topic] {
		go handler(message)
	}
}

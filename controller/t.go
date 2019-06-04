package controller

import (
	"github.com/kgrunwald/goweb/ctx"
	"github.com/kgrunwald/goweb/ilog"
	"github.com/kgrunwald/goweb/pubsub"
)

type T struct {
	logger ilog.Logger
	bus    pubsub.Bus
}

type Message interface {
	GetPayload() string
}

type MessageImpl struct {
	Payload string
}

func (m *MessageImpl) GetPayload() string {
	return m.Payload
}

type AddRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

const Topic = "test topic const"

func (t *T) Add(ctx ctx.Context, a, b int) error {
	res := map[string]int{"result": a + b}
	t.bus.Dispatch(&MessageImpl{"Test Payload"})
	return ctx.OK(res)
}

func (t *T) AddPost(ctx ctx.Context) error {
	req := AddRequest{}
	ctx.Bind(&req)
	res := map[string]int{"result": req.A + req.B}
	return ctx.OK(res)
}

func (t *T) GetVersion() {

}

func (t *T) MessageHandler(msg Message) {
	t.logger.Info("Got message in controller: " + msg.GetPayload())
}

func NewT(logger ilog.Logger, bus pubsub.Bus) *T {
	return &T{logger, bus}
}

package controller

import (
	"net/http"

	"github.com/kgrunwald/goweb/framework"
	"github.com/kgrunwald/goweb/ilog"
	"github.com/kgrunwald/goweb/rest"
)

type T struct {
	logger ilog.Logger
}

type AddRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

func (t *T) Add(r *http.Request, a, b int) framework.Response {
	res := map[string]int{"result": a + b}
	return rest.NewResponse(r, res)
}

func (t *T) AddPost(r *http.Request) framework.Response {
	req := AddRequest{}
	rest.Bind(r, &req)
	res := map[string]int{"result": req.A + req.B}
	return rest.NewResponse(r, res)
}

func NewT(logger ilog.Logger) interface{} {
	return &T{logger}
}

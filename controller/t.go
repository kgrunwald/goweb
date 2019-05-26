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

func (t *T) Add(r *http.Request, a, b int) framework.Response {
	t.logger.Debug(r.URL)
	res := map[string]int{"result": a + b}
	return rest.NewResponse(r, res)
}

func NewT(logger ilog.Logger) interface{} {
	return &T{logger}
}

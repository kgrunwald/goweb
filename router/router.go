package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kgrunwald/goweb/di"
	"github.com/kgrunwald/goweb/ilog"
)

func init() {
	c := di.GetContainer()
	c.Register(NewRouter)
}

type Router struct {
	mux    *mux.Router
	logger ilog.Logger
	ctrls  []interface{}
}

type Route struct {
	Name    string
	Path    string
	Methods []string
}

type RouteHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

func NewRouter(logger ilog.Logger) *Router {
	r := &Router{
		mux:    mux.NewRouter(),
		logger: logger,
	}

	r.mux.Use(LogMiddleware(logger))
	return r
}

func (r *Router) Add(route Route, handler RouteHandler) {
	r.mux.HandleFunc(route.Path, handler.Handle).
		Methods(route.Methods...).
		Name(route.Name)
}

func (r *Router) Start() {
	http.ListenAndServe(":80", r.mux)
}

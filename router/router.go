package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kgrunwald/goweb/di"
	"github.com/kgrunwald/goweb/ilog"
)

func init() {
	c := di.GetContainer()
	c.Register(NewRouter)
}

// Router provides a generic interface for different Routing frameworks.
type Router interface {
	// Add a route to the router
	Add(route Route, handler RouteHandler)

	// PathParams should return any URL parameters from the specified route
	PathParams(req *http.Request) map[string]string

	// Use adds a Middleware handler to the chain of middleware
	Use(fn Middleware)

	// Start listening for incoming connections. This function will block.
	Start(port int)
}

// Route defines a generic route structure
type Route struct {
	Name    string
	Path    string
	Methods []string
}

func (r Route) String() string {
	return fmt.Sprintf("(%s) %s %s", r.Name, r.Methods, r.Path)
}

// Middleware is a function that is invoked before the actual route handler
type Middleware func(next http.Handler) http.Handler

// A RouteHandler is invoked by the Router when a request to the matching Route is received.
type RouteHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

type muxRouter struct {
	mux    *mux.Router
	logger ilog.Logger
	ctrls  []interface{}
}

// NewRouter returns a concrete implementation of the Router interface
func NewRouter(logger ilog.Logger) Router {
	r := &muxRouter{
		mux:    mux.NewRouter(),
		logger: logger,
	}

	r.Use(LogMiddleware(logger))
	return r
}

func (r *muxRouter) Add(route Route, handler RouteHandler) {
	r.mux.HandleFunc(route.Path, handler.Handle).
		Methods(route.Methods...).
		Name(route.Name)
}

func (r *muxRouter) Use(fn Middleware) {
	r.mux.Use((mux.MiddlewareFunc)(fn))
}

func (r *muxRouter) PathParams(req *http.Request) map[string]string {
	return mux.Vars(req)
}

func (r *muxRouter) Start(port int) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), r.mux)
}

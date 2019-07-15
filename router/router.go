package router

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kgrunwald/goweb/ctx"
	"github.com/kgrunwald/goweb/di"
	"github.com/kgrunwald/goweb/ilog"
)

func init() {
	c := di.GetContainer()
	c.Register(NewRouter)
}

// Router provides a generic interface for different Routing frameworks.
type Router interface {
	// Create a new Route object
	NewRoute() Route

	// PathParams should return any URL parameters from the specified route
	PathParams(req *http.Request) map[string]string

	// Use adds a Middleware handler to the chain of middleware
	Use(fn Middleware)

	// Start listening for incoming connections. This function will block.
	Start(port int)
}

// Route defines a generic route structure
type Route interface {
	Name(string) Route
	Path(string) Route
	GetPath() string
	Methods(...string) Route
	Headers(...string) Route
	Handler(f func(http.ResponseWriter, *http.Request)) Route
}

type muxRoute struct {
	route *mux.Route
}

func (r *muxRoute) Handler(f func(http.ResponseWriter, *http.Request)) Route {
	r.route.HandlerFunc(f)
	return r
}

func (r *muxRoute) Name(name string) Route {
	r.route.Name(name)
	return r
}

func (r *muxRoute) GetPath() string {
	path, _ := r.route.GetPathTemplate()
	return path
}

func (r *muxRoute) Path(path string) Route {
	r.route.Path(path)
	return r
}

func (r *muxRoute) Methods(methods ...string) Route {
	r.route.Methods(methods...)
	return r
}

func (r *muxRoute) Headers(headers ...string) Route {
	r.route.Headers(headers...)
	return r
}

// Middleware is a function that is invoked before the actual route handler
type Middleware func(next http.Handler) http.Handler

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

func (r *muxRouter) NewRoute() Route {
	return &muxRoute{
		route: r.mux.NewRoute(),
	}
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

// RouteBinding maps an HTTP route to a controller method to invoke
type RouteBinding struct {
	Route Route
	Vars  []string
}

// RouteHandler implements the HTTP request handler interface by invoking the method specified in the binding.
type RouteHandler struct {
	Method  reflect.Value
	Router  Router
	Binding RouteBinding
	Log     ilog.Logger
}

// Handle is invoked on every incoming HTTP request. It builds up the required parameters for the controller method
// in the `RouteBinding` by inspecting the types of the method arguments. Currently, only primitive
// types may be included. It assumes that each parameter in the controller method correlates
// to a path parameter defined in the `routes.yaml` configuration, and that the parameters are defined in the same order.
// The controller method MUST return an implementation of `Response`.
func (h *RouteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	in := []reflect.Value{}
	method := h.Method.Type()
	numArgs := method.NumIn()

	context := ctx.New(r, w, h.Log)
	if numArgs > 0 {
		in = append(in, reflect.ValueOf(context))

		if len(h.Binding.Vars) > 0 {
			vars := h.Router.PathParams(r)
			for idx, v := range h.Binding.Vars {
				fieldType := h.Method.Type().In(idx + 1).String()
				val, _ := getArgument(vars[v], fieldType)
				in = append(in, val)
			}
		}
	}

	err := h.Method.Call(in)[0].Interface()
	if err != nil {
		context.SendError(err.(error))
	}
}

func getArgument(val, argType string) (reflect.Value, error) {
	switch argType {
	case "string":
		return reflect.ValueOf(val), nil
	case "int":
		v, err := strconv.Atoi(val)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
		return reflect.ValueOf(v), nil
	}

	return reflect.ValueOf(nil), nil
}

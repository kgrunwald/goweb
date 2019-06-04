package framework

import (
	"net/http"
	"reflect"
	"regexp"
	"strconv"

	"github.com/kgrunwald/goweb/ctx"
	"github.com/kgrunwald/goweb/di"
	"github.com/kgrunwald/goweb/ilog"
	"github.com/kgrunwald/goweb/router"
)

// RouteBinding maps an HTTP route to a controller method to invoke
type RouteBinding struct {
	Binding
	Route router.Route
	Vars  []string
}

// RouteHandler implements the HTTP request handler interface by invoking the method specified in the binding.
type RouteHandler struct {
	Method  reflect.Value
	Router  router.Router
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

// InitializeRouter loads the configuration for the router and initializes all of the routes.
func InitializeRouter(r router.Router, logger ilog.Logger, container di.Container) {
	bindings := loadRouteYaml()
	for _, binding := range bindings {
		logger.WithFields(
			"Binding", binding.Binding,
			"Route", binding.Route).
			Debug("Adding route")
		m := container.GetMethod(binding.Service(), binding.Method)
		handler := &RouteHandler{Method: m, Binding: binding, Router: r, Log: logger}
		r.Add(binding.Route, handler)
	}
}

type routesDef struct {
	Routes map[string]struct {
		Path       string
		Methods    []string
		Controller string
		Headers    map[string]string
	} `yaml:",inline"`
}

func loadRouteYaml() []RouteBinding {
	routeFile := routesDef{}
	LoadYaml("routes.yaml", &routeFile)

	bindings := []RouteBinding{}
	for routeName, routeDef := range routeFile.Routes {
		headers := []string{}
		for key, val := range routeDef.Headers {
			headers = append(headers, key, val)
		}

		route := router.Route{
			Name:    routeName,
			Path:    routeDef.Path,
			Methods: routeDef.Methods,
			Headers: headers,
		}

		vars := []string{}
		re := regexp.MustCompile(`\{([^{}]+)\}`)
		matches := re.FindAllStringSubmatch(routeDef.Path, -1)
		for _, match := range matches {
			vars = append(vars, match[1])
		}

		binding := RouteBinding{
			Route:   route,
			Vars:    vars,
			Binding: NewBinding(routeDef.Controller),
		}

		bindings = append(bindings, binding)
	}
	return bindings
}

package framework

import (
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kgrunwald/goweb/di"
	"github.com/kgrunwald/goweb/ilog"
	"github.com/kgrunwald/goweb/router"
)

type Routes struct {
	Routes map[string]struct {
		Path       string
		Methods    []string
		Controller string
	} `yaml:",inline"`
}

type RouteBinding struct {
	Route      router.Route
	Controller string
	Package    string
	Function   string
	Vars       []string
}

type RouteHandler struct {
	Method  reflect.Value
	Binding RouteBinding
	Logger  ilog.Logger
}

type Response interface {
	Send(http.ResponseWriter) error
}

func (h *RouteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	in := []reflect.Value{}
	method := h.Method.Type()
	numArgs := method.NumIn()
	if numArgs > 0 {
		offset := 0
		if method.In(0).String() == "*http.Request" {
			offset = 1
			in = append(in, reflect.ValueOf(r))
		}

		if len(h.Binding.Vars) > 0 {
			vars := mux.Vars(r)
			for idx, v := range h.Binding.Vars {
				fieldType := h.Method.Type().In(idx + offset).String()
				val, _ := getArgument(vars[v], fieldType)
				in = append(in, val)
			}
		}
	}

	out := h.Method.Call(in)
	handler := out[0].Interface().(Response)
	handler.Send(w)
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

func Initialize(r *router.Router, logger ilog.Logger, container di.Container) {
	bindings := LoadRouteYaml()
	ctrls := container.GetControllers()
	for _, binding := range bindings {
		t := reflect.ValueOf(findController(ctrls, binding.Package+"."+binding.Controller))
		m := t.MethodByName(binding.Function)

		handler := &RouteHandler{Method: m, Binding: binding, Logger: logger}
		r.Add(binding.Route, handler)
	}
}

func LoadRouteYaml() []RouteBinding {
	routeFile := Routes{}
	LoadYaml("routes.yaml", &routeFile)

	bindings := []RouteBinding{}
	for routeName, routeDef := range routeFile.Routes {
		parts := strings.Split(routeDef.Controller, "::")
		nameparts := strings.Split(parts[0], ".")
		route := router.Route{
			Name:    routeName,
			Path:    routeDef.Path,
			Methods: routeDef.Methods,
		}

		vars := []string{}
		re := regexp.MustCompile(`\{([^{}]+)\}`)
		matches := re.FindAllStringSubmatch(routeDef.Path, -1)
		for _, match := range matches {
			vars = append(vars, match[1])
		}

		binding := RouteBinding{
			Route:      route,
			Function:   parts[1],
			Package:    nameparts[0],
			Controller: nameparts[1],
			Vars:       vars,
		}

		bindings = append(bindings, binding)
	}
	return bindings
}

func findController(controllers []interface{}, name string) interface{} {
	for _, ctrl := range controllers {
		controllerType := reflect.TypeOf(ctrl).Elem().String()
		if controllerType == name {
			return ctrl
		}
	}

	return nil
}

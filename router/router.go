package router

import (
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kgrunwald/goweb/di"
	"github.com/kgrunwald/goweb/ilog"
	"gopkg.in/yaml.v2"
)

type Route struct {
	Name       string
	Path       string
	Methods    []string
	Controller string
	Package    string
	Function   string
}

type Routes struct {
	Routes map[string]struct {
		Path       string
		Methods    []string
		Controller string
	} `yaml:",inline"`
}

type Router struct {
	mux    *mux.Router
	logger ilog.Logger
	ctrls  []interface{}
}

func NewRouter(logger ilog.Logger, container di.Container) *Router {
	r := &Router{
		mux:    mux.NewRouter(),
		logger: logger,
		ctrls:  container.GetControllers(),
	}

	r.Initialize()

	return r
}

func init() {
	c := di.GetContainer()
	c.Register(NewRouter)
}

func (r *Router) Initialize() {
	routes := LoadYaml()
	for _, route := range routes {
		t := reflect.ValueOf(findController(r.ctrls, route.Package+"."+route.Controller))
		m := t.MethodByName(route.Function)

		d := Dispatcher{Method: m, Path: route.Path}
		// m.Interface().(func(http.ResponseWriter, *http.Request))
		r.mux.HandleFunc(route.Path, d.Dispatch).
			Methods(route.Methods...).
			Name(route.Name)
	}

	r.mux.Use(LogMiddleware(r.logger))
}

func (router *Router) Start() {
	http.ListenAndServe(":80", router.mux)
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

func LoadYaml() []Route {
	path := os.Getenv("CONFIG_DIR") + "/routes.yaml"
	data, _ := ioutil.ReadFile(path)

	routeFile := Routes{}
	yaml.Unmarshal([]byte(data), &routeFile)

	routes := []Route{}
	for routeName, routeDef := range routeFile.Routes {
		parts := strings.Split(routeDef.Controller, "::")
		nameparts := strings.Split(parts[0], ".")
		route := Route{
			Name:       routeName,
			Path:       routeDef.Path,
			Methods:    routeDef.Methods,
			Function:   parts[1],
			Package:    nameparts[0],
			Controller: nameparts[1],
		}

		routes = append(routes, route)
	}
	return routes
}

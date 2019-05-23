package router

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

type Dispatcher struct {
	Method reflect.Value
	Path   string
	Vars   []string
}

//
// func (*T) Add(a, b int) {
// 	vars := mux.Vars(r)
// 	a, _ := strconv.Atoi(vars["a"])
// 	b, _ := strconv.Atoi(vars["b"])
// 	fmt.Fprintf(w, "Result: %d\n", a+b)
// }

func (d *Dispatcher) Dispatch(w http.ResponseWriter, r *http.Request) {
	d.Vars = []string{}
	re := regexp.MustCompile(`\{([^{}]+)\}`)
	matches := re.FindAllStringSubmatch(d.Path, -1)
	for _, match := range matches {
		d.Vars = append(d.Vars, match[1])
	}

	in := []reflect.Value{}
	method := d.Method.Type()
	numArgs := method.NumIn()
	if numArgs > 0 {
		offset := 0
		if method.In(0).String() == "*http.Request" {
			offset = 1
			in = append(in, reflect.ValueOf(r))
		}

		if len(d.Vars) > 0 {
			vars := mux.Vars(r)
			for idx, v := range d.Vars {
				fieldType := d.Method.Type().In(idx + offset).String()
				val, _ := getArgument(vars[v], fieldType)
				in = append(in, val)
			}
		}
	}

	out := d.Method.Call(in)
	fmt.Fprint(w, out[0])
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

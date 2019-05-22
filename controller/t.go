package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type T struct{}

func (*T) Add(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, _ := strconv.Atoi(vars["a"])
	b, _ := strconv.Atoi(vars["b"])
	fmt.Fprintf(w, "Result: %d\n", a+b)
}

func NewT() interface{} {
	return &T{}
}

package controller

import (
	"fmt"
	"net/http"
)

type T struct{}

func (*T) Add(r *http.Request, a, b int) string {
	fmt.Println(r.URL)
	return fmt.Sprintf("Result: %d\n", a+b)
}

func NewT() interface{} {
	return &T{}
}

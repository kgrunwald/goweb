package controller

import (
	"fmt"

	"github.com/kgrunwald/goweb/ilog"
)

type T struct{}

func (*T) Add(a, b int) string {
	return fmt.Sprintf("Result: %d\n", a+b)
}

func NewT(logger ilog.Logger) interface{} {
	logger.Info("Test")
	return &T{}
}

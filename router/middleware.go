package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kgrunwald/goweb/ilog"
)

func LogMiddleware(l ilog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route := mux.CurrentRoute(r)
			l.Info(fmt.Sprintf("Matched route: %s", route.GetName()))
			next.ServeHTTP(w, r)
		})
	}
}

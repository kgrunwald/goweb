package router

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kgrunwald/goweb/ilog"
)

func LogMiddleware(l ilog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			route := mux.CurrentRoute(r)
			l.WithField("Name", route.GetName()).Debug("Matched route")

			sw := &statusWriter{ResponseWriter: w}
			next.ServeHTTP(sw, r)

			duration := time.Now().Sub(start)
			l.WithFields(
				"Duration", duration,
				"Status", sw.status,
				"ContentLen", sw.length,
				"Method", r.Method,
				"RequestURI", r.RequestURI,
			).Info("Access log")
		})
	}
}

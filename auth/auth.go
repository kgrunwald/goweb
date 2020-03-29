package auth

import (
	"os"
	"fmt"
	"net/http"

	"github.com/kgrunwald/goweb"
	"github.com/kgrunwald/goweb/ctx"
	"github.com/kgrunwald/goweb/di"
	"github.com/kgrunwald/goweb/ilog"
	"github.com/kgrunwald/goweb/router"
)

func init() {
	c := di.GetContainer()
	c.Register(NewJWTContext)
	c.Register(NewAPIKeyContext)
}

type Authenticator interface {
	Authenticate(ctx.Context)
	Middleware(next http.Handler) http.Handler
}

type ContextKey int

const (
	ContextKeyAuthenticated ContextKey = iota
	ContextKeyJWTAuthenticated
	ContextKeyClaims
	ContextKeyUserEmail
)

func RequireHTTPS() {
	di.GetContainer().Invoke(func (r router.Router, logger ilog.Logger, info *goweb.ServerInfo) {
		r.Use(HttpsMiddleware(logger, info.Port))
	})
}

const HEADER_FORWARDED_PROTO = "x-forwarded-proto"
const HEADER_HOST = "host"

func HttpsMiddleware(log ilog.Logger, port int) router.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := r.Host
			localhost := fmt.Sprintf("localhost:%d", port)
			headerScheme := r.Header.Get(HEADER_FORWARDED_PROTO)
			if os.Getenv("ENV") != "local" && host != localhost && headerScheme != "https" && r.TLS == nil {
				log.WithFields("host", host, "scheme", headerScheme).Info("Redirecting to HTTPS")
				r.URL.Scheme = "https"
				url := r.URL.String()
				http.Redirect(w, r, url, http.StatusMovedPermanently)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

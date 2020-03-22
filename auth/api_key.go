package auth

import (
	"errors"
	"net/http"
	"os"

	"github.com/kgrunwald/goweb/ctx"
	"github.com/kgrunwald/goweb/ilog"
)

type APIKeyContext struct{
	key string
	log ilog.Logger
}	

type APIKeyScheme struct {
	key string
	log ilog.Logger
}

const (
	ContextKeyApiKeyAuthenticated ContextKey = iota
)

func NewAPIKeyContext(log ilog.Logger) *APIKeyContext {
	key := os.Getenv("API_KEY")
	if len(key) == 0 {
		log.Fatal("$API_KEY environment variable not found")
	}

	return &APIKeyContext{key, log}
}

func (a *APIKeyContext) NewAPIKeyScheme() *APIKeyScheme {
	return &APIKeyScheme{a.key, a.log}
}

func (a *APIKeyScheme) Authenticate(ctx ctx.Context) error {
	header := ctx.Request().Header.Get("x-api-key")
	if header != a.key {
		return errors.New("API key not valid")
	}

	ctx.Log().Info("Authenticated API key")
	ctx.AddValue(ContextKeyAuthenticated, true)
	ctx.AddValue(ContextKeyApiKeyAuthenticated, true)
	return nil
}

func (a *APIKeyScheme) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := ctx.New(r, w, a.log)
		if err := a.Authenticate(context); err != nil {
			context.Forbidden(err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

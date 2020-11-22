package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/kgrunwald/goweb/ctx"
	"github.com/kgrunwald/goweb/ilog"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type JWTContext struct {
	key string
	log ilog.Logger
}

type JWTScheme struct {
	key      string
	log      ilog.Logger
	expected jwt.Expected
}

func NewJWTContext(log ilog.Logger) *JWTContext {
	key := os.Getenv("JWT_SECRET")
	if len(key) == 0 {
		log.Fatal("$JWT_SECRET environment variable not found")
	}

	return &JWTContext{key, log}
}

func (j *JWTContext) UpdateJWTCookie(ctx ctx.Context) error {
	claims := ctx.GetValue(ContextKeyClaims).(*jwt.Claims)
	return j.SetJWTCookie(ctx, claims)
}

func (j *JWTContext) SetJWTCookie(ctx ctx.Context, claims *jwt.Claims) error {
	key := []byte(j.key)
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: key}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return err
	}

	claims.NotBefore = jwt.NewNumericDate(time.Now().UTC())
	raw, err := jwt.Signed(sig).Claims(claims).CompactSerialize()
	if err != nil {
		return err
	}

	cookie := makeCookie(ctx, raw)
	j.log.WithField("cookie", cookie.String()).Debug("Setting cookie")
	http.SetCookie(ctx.Writer(), cookie)
	return nil
}

func makeCookie(ctx ctx.Context, token string) *http.Cookie {
	proto := ctx.Request().Header.Get("X-Forwarded-Proto")
	secure := proto == "https"
	return &http.Cookie{
		Name:     "authorization",
		Path:     "/",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   secure,
	}
}

func (j *JWTContext) DeleteJWTCookie(ctx ctx.Context) {
	cookie := makeCookie(ctx, "")
	cookie.Expires = time.Now()
	http.SetCookie(ctx.Writer(), cookie)
}

func (j *JWTContext) NewJWTScheme(issuer, audience string) *JWTScheme {
	expected := jwt.Expected{
		Issuer:   issuer,
		Audience: jwt.Audience{audience},
	}
	return &JWTScheme{
		key:      j.key,
		log:      j.log,
		expected: expected,
	}
}

func (j *JWTScheme) Authenticate(ctx ctx.Context) error {
	authToken, err := ctx.Request().Cookie("authorization")
	log := ctx.Log()
	if err != nil || len(authToken.Value) == 0 {
		log.Info("No JWT cookie provided")
		return errors.New("No JWT cookie provided")
	}

	tok, err := jwt.ParseSigned(string(authToken.Value))
	if err != nil {
		log.WithField("error", err).Error("Could not parse JWT")
		return err
	}

	claims := jwt.Claims{}
	if err := tok.Claims([]byte(j.key), &claims); err != nil {
		log.WithField("error", err).Error("Could not validate signature")
		return err
	}

	if err := claims.Validate(j.expected); err != nil {
		log.WithFields("claims", claims, "expected", j.expected).Error("Claims did not match expected")
	}

	ctx.Log().WithField("claims", claims).Info("Authenticated user")
	ctx.AddValue(ContextKeyAuthenticated, true)
	ctx.AddValue(ContextKeyJWTAuthenticated, true)
	ctx.AddValue(ContextKeyClaims, &claims)
	ctx.AddValue(ContextKeyUserEmail, claims.Subject)
	return nil
}

func (j *JWTScheme) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := ctx.New(r, w, j.log)
		if err := j.Authenticate(context); err != nil {
			context.Forbidden(err)
			return
		}
		next.ServeHTTP(context.Writer(), context.Request())
	})
}

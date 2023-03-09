package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
	"github.com/hiennguyen9874/stockk-go/pkg/jwt"
)

var (
	TokenCtxKey = &contextKey{"Token"}
	IdCtxKey    = &contextKey{"Id"}
	EmailCtxKey = &contextKey{"Email"}
	ErrorCtxKey = &contextKey{"Error"}
	UserCtxKey  = &contextKey{"User"}
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "jwtauth context value " + k.name
}

// Verifier http middleware handler will verify a JWT string from a http request.
//
// Verifier will search for a JWT token in a http request, in the order:
//  1. 'Authorization: BEARER T' request header
//  2. Cookie 'jwt' value
//
// The first JWT string that is found as a query parameter, authorization header
// or cookie header is then decoded by the `jwt-go` library and a *jwt.Token
// object is set on the request context. In the case of a signature decoding error
// the Verifier will also set the error on the request context.
//
// The Verifier always calls the next http handler in sequence, which can either
// be the generic `jwtauth.Authenticator` middleware or your own custom handler
// which checks the request context jwt token and error to prepare a custom
// http response.
func (mw *MiddlewareManager) Verifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token := TokenFromHeader(r)

		if token == "" {
			err := httpErrors.ErrorTokenNotFound
			ctx = context.WithValue(ctx, ErrorCtxKey, err)
		} else {
			id, email, err := jwt.ParseToken(token, mw.cfg.Server.JwtSecretKey)
			ctx = context.WithValue(ctx, TokenCtxKey, token)
			ctx = context.WithValue(ctx, IdCtxKey, id)
			ctx = context.WithValue(ctx, EmailCtxKey, email)
			ctx = context.WithValue(ctx, ErrorCtxKey, err)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Authenticator is a default authentication middleware to enforce access from the
// Verifier middleware request context values. The Authenticator sends a 401 Unauthorized
// response for any unverified tokens and passes the good ones through. It's just fine
// until you decide to write something similar and customize your client response.
func (mw *MiddlewareManager) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err, _ := r.Context().Value(ErrorCtxKey).(error)

		if err != nil {
			render.Render(w, r, httpErrors.ErrRender(err))
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

func (mw *MiddlewareManager) CurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, _ := r.Context().Value(IdCtxKey).(string)
		err, _ := r.Context().Value(ErrorCtxKey).(error)

		if err != nil || id == "" {
			render.Render(w, r, httpErrors.ErrRender(err))
			return
		}

		idParsed, err := uuid.Parse(id)

		if err != nil {
			render.Render(w, r, httpErrors.ErrRender(httpErrors.Err(err, http.StatusBadRequest, httpErrors.ErrorInvalidJWTClaims.Error())))
			return
		}

		user, err := mw.usersUC.Get(ctx, idParsed)

		if err != nil {
			render.Render(w, r, httpErrors.ErrRender(err))
			return
		}

		ctx = context.WithValue(ctx, UserCtxKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mw *MiddlewareManager) SuperUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := GetUserFromCtx(ctx)

		if err != nil {
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ParseErrors(err)))
			return
		}

		if !mw.usersUC.IsSuper(ctx, *user) {
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ErrNotEnoughPrivileges(errors.New("not enough privileges"))))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (mw *MiddlewareManager) ActiveUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := GetUserFromCtx(ctx)

		if err != nil {
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ParseErrors(err)))
			return
		}

		if !mw.usersUC.IsActive(ctx, *user) {
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ErrInactiveUser(errors.New("inactive user"))))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// TokenFromHeader tries to retreive the token string from the
// "Authorization" reqeust header: "Authorization: BEARER T".
func TokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

// Get user from context
func GetUserFromCtx(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value(UserCtxKey).(*models.User)
	if !ok {
		return nil, httpErrors.ErrUnauthorized(errors.New("Can convert user from context"))
	}
	return user, nil
}

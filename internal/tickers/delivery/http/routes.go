package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/hiennguyen9874/go-boilerplate/internal/middleware"
	"github.com/hiennguyen9874/go-boilerplate/internal/tickers"
)

func MapTickerRoute(router *chi.Mux, h tickers.Handlers, mw *middleware.MiddlewareManager) {
	// User routes
	router.Route("/ticker", func(r chi.Router) {
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Get("/", h.GetMulti())
			// Per id routes
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.Get())
				r.Delete("/", h.Delete())
			})
		})
	})
}

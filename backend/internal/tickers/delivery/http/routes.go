package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/hiennguyen9874/stockk-go/internal/middleware"
	"github.com/hiennguyen9874/stockk-go/internal/tickers"
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
			// Per symbol routes
			r.Route("/{symbol}", func(r chi.Router) {
				r.Get("/", h.GetBySymbol())
				r.Patch("/", h.UpdateIsActiveBySymbol())
			})
		})
	})
}

package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/hiennguyen9874/stockk-go/internal/middleware"
	"github.com/hiennguyen9874/stockk-go/internal/stockssnapshot"
)

func MapStockSnapshotRoute(router *chi.Mux, h stockssnapshot.Handlers, mw *middleware.MiddlewareManager) {
	// User routes
	router.Route("/stocksnapshot", func(r chi.Router) {
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			// Per symbol routes
			r.Route("/{symbol}", func(r chi.Router) {
				r.Get("/", h.GetStockSnapshotBySymbol())
			})
		})
	})
}

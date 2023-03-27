package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/hiennguyen9874/stockk-go/internal/middleware"
	"github.com/hiennguyen9874/stockk-go/internal/watchlists"
)

func MapWatchListRoute(router *chi.Mux, h watchlists.Handlers, mw *middleware.MiddlewareManager) {
	// WatchList routes
	router.Route("/watchlist", func(r chi.Router) {
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Get("/", h.GetMulti())
			r.Post("/", h.Create())
			// Per id routes
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.Get())
				// Admin routes
				r.Delete("/", h.Delete())
				r.Put("/", h.Update())
			})
		})
	})
}

package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/hiennguyen9874/stockk-go/internal/clients"
	"github.com/hiennguyen9874/stockk-go/internal/middleware"
)

func MapClientRoute(router *chi.Mux, h clients.Handlers, mw *middleware.MiddlewareManager) {
	// Client routes
	router.Route("/client", func(r chi.Router) {
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

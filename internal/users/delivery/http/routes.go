package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/hiennguyen9874/stockk-go/internal/middleware"
	"github.com/hiennguyen9874/stockk-go/internal/users"
	"gorm.io/gorm"
)

func MapUserRoute(router *chi.Mux, db *gorm.DB, h users.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/user", func(r chi.Router) {
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier)
			r.Use(mw.Authenticator)
			r.Use(mw.CurrentUser)
			r.Use(mw.ActiveUser)
			r.Get("/me", h.Me())
			r.Post("/me", h.UpdateMe())
			r.Post("/me/pass", h.UpdatePasswordMe())
			// Admin routes
			r.Group(func(r chi.Router) {
				r.Use(mw.SuperUser)
				r.Get("/", h.GetMulti())
				r.Post("/", h.Create())
			})
			// Per id routes
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.Get())
				// Admin routes
				r.Group(func(r chi.Router) {
					r.Use(mw.SuperUser)
					r.Delete("/", h.Delete())
					r.Post("/", h.Update())
					r.Post("/pass", h.UpdatePassword())
				})
			})
		})
		// Public routes
		r.Group(func(r chi.Router) {
			r.Post("/login", h.SignIn())
		})
	})
}

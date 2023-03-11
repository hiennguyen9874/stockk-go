package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/hiennguyen9874/stockk-go/internal/middleware"
	"github.com/hiennguyen9874/stockk-go/internal/users"
)

func MapUserRoute(router *chi.Mux, h users.Handlers, mw *middleware.MiddlewareManager) {
	// Auth routes
	router.Route("/auth", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Post("/login", h.SignIn())
			r.Get("/publickey", h.GetPublicKey())
		})
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(false))
			r.Use(mw.Authenticator())
			r.Get("/refresh", h.RefreshToken())
			r.Get("/logout", h.Logout())
			r.Get("/logoutall", h.LogoutAllToken())
		})
	})
	// User routes
	router.Route("/user", func(r chi.Router) {
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Get("/me", h.Me())
			r.Put("/me", h.UpdateMe())
			r.Patch("/me/pass", h.UpdatePasswordMe())
			// Admin routes
			r.Group(func(r chi.Router) {
				r.Use(mw.SuperUser())
				r.Get("/", h.GetMulti())
				r.Post("/", h.Create())
			})
			// Per id routes
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.Get())
				// Admin routes
				r.Group(func(r chi.Router) {
					r.Use(mw.SuperUser())
					r.Delete("/", h.Delete())
					r.Put("/", h.Update())
					r.Patch("/pass", h.UpdatePassword())
					r.Get("/logoutall", h.LogoutAllAdmin())
				})
			})
		})
	})
}

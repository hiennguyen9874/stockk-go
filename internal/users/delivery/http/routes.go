package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/hiennguyen9874/stockk-go/internal/users"
	"gorm.io/gorm"
)

func MapUserRoute(router *chi.Mux, db *gorm.DB, h users.Handlers) {
	router.Route("/user", func(r chi.Router) {
		r.Get("/", h.GetMulti())
		r.Post("/", h.Create())
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.Get())
			r.Delete("/", h.Delete())
		})
	})
}

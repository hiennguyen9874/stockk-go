package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/hiennguyen9874/stockk-go/internal/drawingtemplates"
	"github.com/hiennguyen9874/stockk-go/internal/middleware"
)

func MapDrawingTemplateRoute(router *chi.Mux, h drawingtemplates.Handlers, mw *middleware.MiddlewareManager) {
	// User routes
	router.Route("/drawing_templates", func(r chi.Router) {
		r.Get("/", h.Get())
		r.Post("/", h.CreateOrUpdate())
		r.Delete("/", h.Delete())
	})
}

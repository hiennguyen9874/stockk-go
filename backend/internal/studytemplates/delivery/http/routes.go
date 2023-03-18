package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/hiennguyen9874/stockk-go/internal/middleware"
	"github.com/hiennguyen9874/stockk-go/internal/studytemplates"
)

func MapStudyTemplateRoute(router *chi.Mux, h studytemplates.Handlers, mw *middleware.MiddlewareManager) {
	// User routes
	router.Route("/study_templates", func(r chi.Router) {
		r.Get("/", h.Get())
		r.Post("/", h.CreateOrUpdate())
		r.Delete("/", h.Delete())
	})
}

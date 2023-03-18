package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/hiennguyen9874/stockk-go/internal/charts"
	"github.com/hiennguyen9874/stockk-go/internal/middleware"
)

func MapChartRoute(router *chi.Mux, h charts.Handlers, mw *middleware.MiddlewareManager) {
	// User routes
	router.Route("/charts", func(r chi.Router) {
		r.Get("/", h.Get())
		r.Post("/", h.CreateOrUpdate())
		r.Delete("/", h.Delete())
	})
}

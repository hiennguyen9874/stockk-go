package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/hiennguyen9874/stockk-go/internal/dchart"
	"github.com/hiennguyen9874/stockk-go/internal/middleware"
)

func MapDchartRoute(router *chi.Mux, h dchart.Handlers, mw *middleware.MiddlewareManager) {
	router.Route("/dchart", func(r chi.Router) {
		r.Get("/time", h.GetTime())
		r.Get("/config", h.GetConfig())
		r.Get("/symbols", h.GetSymbols())
		r.Get("/search", h.Search())
		r.Get("/history", h.History())
	})
}

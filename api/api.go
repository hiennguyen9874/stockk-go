package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/hiennguyen9874/stockk-go/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(enableCORS bool) (*chi.Mux, error) {
	db, err := gorm.Open(postgres.Open(config.GetDNSConfig()), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	RegisterRoutes(r, db)

	return r, nil
}

func RegisterRoutes(router *chi.Mux, db *gorm.DB) {
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	UserRoute(router, db)
}

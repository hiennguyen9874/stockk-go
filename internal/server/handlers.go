package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"gorm.io/gorm"

	"github.com/hiennguyen9874/stockk-go/config"
	apiMiddleware "github.com/hiennguyen9874/stockk-go/internal/middleware"
	userHttp "github.com/hiennguyen9874/stockk-go/internal/users/delivery/http"
	userRepository "github.com/hiennguyen9874/stockk-go/internal/users/repository"
	userUseCase "github.com/hiennguyen9874/stockk-go/internal/users/usecase"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

func New(db *gorm.DB, cfg *config.Config, logger logger.Logger) (*chi.Mux, error) {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(apiMiddleware.Cors(cfg)))

	RegisterRoutes(r, db, cfg, logger)

	return r, nil
}

func RegisterRoutes(router *chi.Mux, db *gorm.DB, cfg *config.Config, logger logger.Logger) {
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// routes.UserRoute(router, db)

	// Repository
	userRepo := userRepository.CreateUserRepository(db)

	// UseCase
	userUC := userUseCase.CreateUserUseCaseI(userRepo, cfg, logger)

	// Handler
	userHandler := userHttp.CreateUserHandler(userUC, cfg, logger)

	// middleware
	mw := apiMiddleware.CreateMiddlewareManager(cfg, logger, userUC)

	userHttp.MapUserRoute(router, db, userHandler, mw)
}

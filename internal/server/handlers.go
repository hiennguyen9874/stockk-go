package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/hiennguyen9874/stockk-go/config"
	apiMiddleware "github.com/hiennguyen9874/stockk-go/internal/middleware"
	tickerHttp "github.com/hiennguyen9874/stockk-go/internal/tickers/delivery/http"
	tickerRepository "github.com/hiennguyen9874/stockk-go/internal/tickers/repository"
	tickerUseCase "github.com/hiennguyen9874/stockk-go/internal/tickers/usecase"
	userHttp "github.com/hiennguyen9874/stockk-go/internal/users/delivery/http"
	userRepository "github.com/hiennguyen9874/stockk-go/internal/users/repository"
	userUseCase "github.com/hiennguyen9874/stockk-go/internal/users/usecase"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

func New(db *gorm.DB, redisClient *redis.Client, cfg *config.Config, logger logger.Logger) (*chi.Mux, error) {
	r := chi.NewRouter()

	// Repository
	userPgRepo := userRepository.CreateUserPgRepository(db)
	userRedisRepo := userRepository.CreateUserRedisRepository(redisClient)

	tickerPgRepo := tickerRepository.CreateTickerPgRepository(db)

	// UseCase
	userUC := userUseCase.CreateUserUseCaseI(userPgRepo, userRedisRepo, cfg, logger)
	tickerUC := tickerUseCase.CreateTickerUseCaseI(tickerPgRepo, cfg, logger)

	// Handler
	userHandler := userHttp.CreateUserHandler(userUC, cfg, logger)
	tickerHandler := tickerHttp.CreateTickerHandler(tickerUC, cfg, logger)

	// middleware
	mw := apiMiddleware.CreateMiddlewareManager(cfg, logger, userUC)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(time.Duration(cfg.Server.ProcessTimeout) * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(mw.Cors()))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	userHttp.MapUserRoute(r, userHandler, mw)
	tickerHttp.MapTickerRoute(r, tickerHandler, mw)

	return r, nil
}

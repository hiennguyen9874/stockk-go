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
	authHttp "github.com/hiennguyen9874/stockk-go/internal/auth/delivery/http"
	barRepository "github.com/hiennguyen9874/stockk-go/internal/bars/repository"
	barUseCase "github.com/hiennguyen9874/stockk-go/internal/bars/usecase"
	chartHttp "github.com/hiennguyen9874/stockk-go/internal/charts/delivery/http"
	chartRepository "github.com/hiennguyen9874/stockk-go/internal/charts/repository"
	chartUseCase "github.com/hiennguyen9874/stockk-go/internal/charts/usecase"
	dchartHttp "github.com/hiennguyen9874/stockk-go/internal/dchart/delivery/http"
	apiMiddleware "github.com/hiennguyen9874/stockk-go/internal/middleware"
	studyTemplateHttp "github.com/hiennguyen9874/stockk-go/internal/studytemplates/delivery/http"
	studyTemplateRepository "github.com/hiennguyen9874/stockk-go/internal/studytemplates/repository"
	studyTemplateUseCase "github.com/hiennguyen9874/stockk-go/internal/studytemplates/usecase"
	tickerHttp "github.com/hiennguyen9874/stockk-go/internal/tickers/delivery/http"
	tickerRepository "github.com/hiennguyen9874/stockk-go/internal/tickers/repository"
	tickerUseCase "github.com/hiennguyen9874/stockk-go/internal/tickers/usecase"
	userHttp "github.com/hiennguyen9874/stockk-go/internal/users/delivery/http"
	userRepository "github.com/hiennguyen9874/stockk-go/internal/users/repository"
	userUseCase "github.com/hiennguyen9874/stockk-go/internal/users/usecase"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func New(db *gorm.DB, redisClient *redis.Client, influxDB influxdb2.Client, cfg *config.Config, logger logger.Logger) (*chi.Mux, error) {
	r := chi.NewRouter()

	// Repository
	userPgRepo := userRepository.CreateUserPgRepository(db)
	userRedisRepo := userRepository.CreateUserRedisRepository(redisClient)
	tickerPgRepo := tickerRepository.CreateTickerPgRepository(db)
	barInfluxDBRepo := barRepository.CreateBarRepo(influxDB, "history")
	barRedisRepo := barRepository.CreateBarRedisRepository(redisClient)
	chartPgRepo := chartRepository.CreateChartPgRepository(db)
	studyTemplatePgRepo := studyTemplateRepository.CreateStudyTemplatePgRepository(db)

	// UseCase
	userUC := userUseCase.CreateUserUseCaseI(userPgRepo, userRedisRepo, cfg, logger)
	tickerUC := tickerUseCase.CreateTickerUseCaseI(tickerPgRepo, cfg, logger)
	barUseCase := barUseCase.CreateBarUseCaseI(barInfluxDBRepo, barRedisRepo, tickerPgRepo, cfg, logger)
	chartUseCase := chartUseCase.CreateChartUseCaseI(chartPgRepo, cfg, logger)
	studyTemplateUseCase := studyTemplateUseCase.StudyTemplateUseCaseI(studyTemplatePgRepo, cfg, logger)

	// Handler
	userHandler := userHttp.CreateUserHandler(userUC, cfg, logger)
	authHandler := authHttp.CreateAuthHandler(userUC, cfg, logger)
	tickerHandler := tickerHttp.CreateTickerHandler(tickerUC, cfg, logger)
	dchartHandler := dchartHttp.CreateDchartHandler(tickerUC, barUseCase, cfg, logger)
	chartHandler := chartHttp.CreateChartHandler(chartUseCase, cfg, logger)
	studyTemplateHandler := studyTemplateHttp.CreateStudyTemplateHandler(studyTemplateUseCase, cfg, logger)

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

	apiRouter := chi.NewRouter()
	r.Mount("/api", apiRouter)

	apiRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	authHttp.MapAuthRoute(apiRouter, authHandler, mw)
	userHttp.MapUserRoute(apiRouter, userHandler, mw)
	tickerHttp.MapTickerRoute(apiRouter, tickerHandler, mw)
	dchartHttp.MapDchartRoute(apiRouter, dchartHandler, mw)

	// Storage api
	storageRouter := chi.NewRouter()
	apiRouter.Mount("/storage/1.1", storageRouter)

	chartHttp.MapChartRoute(storageRouter, chartHandler, mw)
	studyTemplateHttp.MapStudyTemplatesRoute(storageRouter, studyTemplateHandler, mw)

	return r, nil
}

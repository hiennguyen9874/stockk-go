package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/hibiken/asynq"
	_ "github.com/hiennguyen9874/stockk-go/docs" // docs is generated by Swag CLI, you have to import it.
	httpSwagger "github.com/hiennguyen9874/stockk-go/pkg/http-swagger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/hiennguyen9874/stockk-go/config"
	authHttp "github.com/hiennguyen9874/stockk-go/internal/auth/delivery/http"
	barRepository "github.com/hiennguyen9874/stockk-go/internal/bars/repository"
	barUseCase "github.com/hiennguyen9874/stockk-go/internal/bars/usecase"
	chartHttp "github.com/hiennguyen9874/stockk-go/internal/charts/delivery/http"
	chartRepository "github.com/hiennguyen9874/stockk-go/internal/charts/repository"
	chartUseCase "github.com/hiennguyen9874/stockk-go/internal/charts/usecase"
	clientHttp "github.com/hiennguyen9874/stockk-go/internal/clients/delivery/http"
	clientRepository "github.com/hiennguyen9874/stockk-go/internal/clients/repository"
	clientUseCase "github.com/hiennguyen9874/stockk-go/internal/clients/usecase"
	dchartHttp "github.com/hiennguyen9874/stockk-go/internal/dchart/delivery/http"
	drawingTemplateHttp "github.com/hiennguyen9874/stockk-go/internal/drawingtemplates/delivery/http"
	drawingTemplateRepository "github.com/hiennguyen9874/stockk-go/internal/drawingtemplates/repository"
	drawingTemplateUseCase "github.com/hiennguyen9874/stockk-go/internal/drawingtemplates/usecase"
	apiMiddleware "github.com/hiennguyen9874/stockk-go/internal/middleware"
	stockSnapshotHttp "github.com/hiennguyen9874/stockk-go/internal/stockssnapshot/delivery/http"
	stockSnapshotRepository "github.com/hiennguyen9874/stockk-go/internal/stockssnapshot/repository"
	stockSnapshotUseCase "github.com/hiennguyen9874/stockk-go/internal/stockssnapshot/usecase"
	studyTemplateHttp "github.com/hiennguyen9874/stockk-go/internal/studytemplates/delivery/http"
	studyTemplateRepository "github.com/hiennguyen9874/stockk-go/internal/studytemplates/repository"
	studyTemplateUseCase "github.com/hiennguyen9874/stockk-go/internal/studytemplates/usecase"
	tickerHttp "github.com/hiennguyen9874/stockk-go/internal/tickers/delivery/http"
	tickerRepository "github.com/hiennguyen9874/stockk-go/internal/tickers/repository"
	tickerUseCase "github.com/hiennguyen9874/stockk-go/internal/tickers/usecase"
	userHttp "github.com/hiennguyen9874/stockk-go/internal/users/delivery/http"
	userDistributor "github.com/hiennguyen9874/stockk-go/internal/users/distributor"
	userRepository "github.com/hiennguyen9874/stockk-go/internal/users/repository"
	userUseCase "github.com/hiennguyen9874/stockk-go/internal/users/usecase"
	watchListHttp "github.com/hiennguyen9874/stockk-go/internal/watchlists/delivery/http"
	watchListRepository "github.com/hiennguyen9874/stockk-go/internal/watchlists/repository"
	watchListUseCase "github.com/hiennguyen9874/stockk-go/internal/watchlists/usecase"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// @title Stockk Go
// @version 1.0

// @BasePath /api
// @securitydefinitions.oauth2.password	OAuth2Password
// @tokenUrl /api/auth/login
func New(db *gorm.DB, redisClient *redis.Client, taskRedisClient *asynq.Client, influxDB influxdb2.Client, cfg *config.Config, logger logger.Logger) (*chi.Mux, error) {
	r := chi.NewRouter()

	// Repository
	userPgRepo := userRepository.CreateUserPgRepository(db)
	userRedisRepo := userRepository.CreateUserRedisRepository(redisClient)
	tickerPgRepo := tickerRepository.CreateTickerPgRepository(db)
	tickerRedisRepo := tickerRepository.CreateTickerRedisRepository(redisClient)
	barInfluxDBRepo := barRepository.CreateBarRepo(influxDB, cfg.InfluxDB.Org)
	barRedisRepo := barRepository.CreateBarRedisRepository(redisClient)
	chartPgRepo := chartRepository.CreateChartPgRepository(db)
	studyTemplatePgRepo := studyTemplateRepository.CreateStudyTemplatePgRepository(db)
	drawingTemplatePgRepo := drawingTemplateRepository.CreateDrawingTemplatePgRepository(db)
	clientPgRepo := clientRepository.CreateClientPgRepository(db)
	watchListPgRePo := watchListRepository.CreateWatchListPgRepository(db)
	stockSnapshotRedisRepo := stockSnapshotRepository.CreateStockSnapshotRedisRepository(redisClient)

	// Distributor
	userRedisTaskDistributor := userDistributor.NewUserRedisTaskDistributor(taskRedisClient, cfg, logger)

	// UseCase
	userUC := userUseCase.CreateUserUseCaseI(userPgRepo, userRedisRepo, userRedisTaskDistributor, cfg, logger)
	tickerUC := tickerUseCase.CreateTickerUseCaseI(tickerPgRepo, tickerRedisRepo, cfg, logger)
	barUseCase := barUseCase.CreateBarUseCaseI(barInfluxDBRepo, barRedisRepo, tickerPgRepo, tickerRedisRepo, cfg, logger)
	chartUseCase := chartUseCase.CreateChartUseCaseI(chartPgRepo, cfg, logger)
	studyTemplateUseCase := studyTemplateUseCase.StudyTemplateUseCaseI(studyTemplatePgRepo, cfg, logger)
	drawingTemplateUseCase := drawingTemplateUseCase.DrawingTemplateUseCaseI(drawingTemplatePgRepo, cfg, logger)
	clientUC := clientUseCase.CreateClientUseCaseI(clientPgRepo, cfg, logger)
	watchListUC := watchListUseCase.CreateWatchListUseCaseI(watchListPgRePo, cfg, logger)
	stockSnapshotUC := stockSnapshotUseCase.CreateTickerUseCaseI(tickerUC, stockSnapshotRedisRepo, cfg, logger)

	// Handler
	userHandler := userHttp.CreateUserHandler(userUC, cfg, logger)
	authHandler := authHttp.CreateAuthHandler(userUC, cfg, logger)
	tickerHandler := tickerHttp.CreateTickerHandler(tickerUC, cfg, logger)
	dchartHandler := dchartHttp.CreateDchartHandler(tickerUC, barUseCase, cfg, logger)
	chartHandler := chartHttp.CreateChartHandler(chartUseCase, cfg, logger)
	studyTemplateHandler := studyTemplateHttp.CreateStudyTemplateHandler(studyTemplateUseCase, cfg, logger)
	drawingTemplateHandler := drawingTemplateHttp.CreateDrawingTemplateHandler(drawingTemplateUseCase, cfg, logger)
	clientHandler := clientHttp.CreateClientHandler(clientUC, cfg, logger)
	watchListHandler := watchListHttp.CreateWatchListHandler(watchListUC, cfg, logger)
	stockSnapshothandler := stockSnapshotHttp.CreatestockSnapshotHandler(stockSnapshotUC, cfg, logger)

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

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
	))

	apiRouter := chi.NewRouter()
	r.Mount("/api", apiRouter)

	apiRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, "pong")
	})

	authHttp.MapAuthRoute(apiRouter, authHandler, mw)
	userHttp.MapUserRoute(apiRouter, userHandler, mw)
	tickerHttp.MapTickerRoute(apiRouter, tickerHandler, mw)
	dchartHttp.MapDchartRoute(apiRouter, dchartHandler, mw)
	clientHttp.MapClientRoute(apiRouter, clientHandler, mw)
	watchListHttp.MapWatchListRoute(apiRouter, watchListHandler, mw)
	stockSnapshotHttp.MapStockSnapshotRoute(apiRouter, stockSnapshothandler, mw)

	// Storage api
	storageRouter := chi.NewRouter()
	apiRouter.Mount("/storage/1.1", storageRouter)

	chartHttp.MapChartRoute(storageRouter, chartHandler, mw)
	studyTemplateHttp.MapStudyTemplateRoute(storageRouter, studyTemplateHandler, mw)
	drawingTemplateHttp.MapDrawingTemplateRoute(storageRouter, drawingTemplateHandler, mw)

	return r, nil
}

package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hiennguyen9874/stockk-go/config"
	barRepository "github.com/hiennguyen9874/stockk-go/internal/bars/repository"
	barUseCase "github.com/hiennguyen9874/stockk-go/internal/bars/usecase"
	tickerRepository "github.com/hiennguyen9874/stockk-go/internal/tickers/repository"
	"github.com/hiennguyen9874/stockk-go/pkg/db/influxdb"
	"github.com/hiennguyen9874/stockk-go/pkg/db/postgres"
	"github.com/hiennguyen9874/stockk-go/pkg/db/redis"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/hiennguyen9874/stockk-go/pkg/sentry"
	"github.com/spf13/cobra"
)

var crawlHistoryDCmd = &cobra.Command{
	Use:   "crawlhistoryd",
	Short: "crawl history",
	Long:  "crawl history",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		cfg := config.GetCfg()

		appLogger := logger.NewApiLogger(cfg)
		appLogger.InitLogger()
		appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

		sentry.Init(cfg)
		defer sentry.Flush()

		psqlDB, err := postgres.NewPsqlDB(cfg)
		if err != nil {
			appLogger.Fatalf("Postgresql init: %s", err)
		} else {
			appLogger.Infof("Postgres connected")
		}

		influxDB, err := influxdb.NewInfluxDB(cfg)
		if err != nil {
			appLogger.Fatalf("InfluxDB init: %s", err)
		} else {
			appLogger.Infof("InfluxDB connected")
		}

		redisClient := redis.NewRedis(cfg)

		// Repository
		tickerPgRepo := tickerRepository.CreateTickerPgRepository(psqlDB)
		tickerRedisRepo := tickerRepository.CreateTickerRedisRepository(redisClient)
		barInfluxDBRepo := barRepository.CreateBarRepo(influxDB, cfg.InfluxDB.Org)
		barRedisRepo := barRepository.CreateBarRedisRepository(redisClient)

		barUseCase := barUseCase.CreateBarUseCaseI(barInfluxDBRepo, barRedisRepo, tickerPgRepo, tickerRedisRepo, cfg, appLogger)

		go func() {
			for {
				status, err := influxDB.Ping(ctx)
				if err != nil {
					appLogger.Warn(err)
					time.Sleep(1 * time.Hour)
					continue
				}
				if !status {
					appLogger.Warn("influxdb not connected")
					time.Sleep(1 * time.Hour)
					continue
				}

				appLogger.Info("Start syncing....")
				err = barUseCase.SyncDAllSymbol(ctx, cfg.Crawler.TickerDownloadBatchSize, cfg.Crawler.TickerInsertBatchSize, cfg.Crawler.BarInsertBatchSize)
				if err != nil {
					appLogger.Warn(err)
				}
				appLogger.Info("Done sync, sleep 30s!")

				time.Sleep(1 * time.Hour)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		sig := <-quit

		appLogger.Infof("Shutting down server... Reason: %s", sig)
	},
}

func init() {
	RootCmd.AddCommand(crawlHistoryDCmd)
}

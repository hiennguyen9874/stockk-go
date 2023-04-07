package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hiennguyen9874/stockk-go/config"
	stockSnapshotRepository "github.com/hiennguyen9874/stockk-go/internal/stockssnapshot/repository"
	stockSnapshotUseCase "github.com/hiennguyen9874/stockk-go/internal/stockssnapshot/usecase"
	tickerRepository "github.com/hiennguyen9874/stockk-go/internal/tickers/repository"
	tickerUseCase "github.com/hiennguyen9874/stockk-go/internal/tickers/usecase"
	"github.com/hiennguyen9874/stockk-go/pkg/db/postgres"
	"github.com/hiennguyen9874/stockk-go/pkg/db/redis"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/hiennguyen9874/stockk-go/pkg/sentry"
	"github.com/spf13/cobra"
)

var stockSnapshotCmd = &cobra.Command{
	Use:   "crawlstocksnapshot",
	Short: "crawl stock snapshot",
	Long:  "crawl stock snapshot",
	Run: func(cmd *cobra.Command, args []string) {
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

		redisClient := redis.NewRedis(cfg)

		// Repository
		tickerPgRepo := tickerRepository.CreateTickerPgRepository(psqlDB)
		tickerRedisRepo := tickerRepository.CreateTickerRedisRepository(redisClient)
		stockSnapshotRedisRepo := stockSnapshotRepository.CreateStockSnapshotRedisRepository(redisClient)

		// UseCase
		tickerUC := tickerUseCase.CreateTickerUseCaseI(tickerPgRepo, tickerRedisRepo, cfg, appLogger)
		stockSnapshotUC := stockSnapshotUseCase.CreateTickerUseCaseI(tickerUC, stockSnapshotRedisRepo, cfg, appLogger)

		ctx := context.Background()

		go func() {
			for {
				stockSnapshotUC.CrawlAllStocksSnapshot(ctx)

				time.Sleep(30 * time.Second)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		sig := <-quit

		appLogger.Infof("Shutting down server... Reason: %s", sig)
	},
}

func init() {
	RootCmd.AddCommand(stockSnapshotCmd)
}

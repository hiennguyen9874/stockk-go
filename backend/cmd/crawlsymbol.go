package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hiennguyen9874/stockk-go/config"
	tickerRepository "github.com/hiennguyen9874/stockk-go/internal/tickers/repository"
	tickerUseCase "github.com/hiennguyen9874/stockk-go/internal/tickers/usecase"
	"github.com/hiennguyen9874/stockk-go/pkg/db/postgres"
	"github.com/hiennguyen9874/stockk-go/pkg/db/redis"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/hiennguyen9874/stockk-go/pkg/sentry"
	"github.com/spf13/cobra"
)

var crawlSymbolCmd = &cobra.Command{
	Use:   "crawlsymbol",
	Short: "crawl symbol",
	Long:  "crawl symbol",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetCfg()

		appLogger := logger.NewApiLogger(cfg)
		appLogger.InitLogger()
		appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

		err := sentry.Init(cfg)
		if err != nil {
			appLogger.Fatalf("Sentry init: %s", err)
		}
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

		// UseCase
		tickerUC := tickerUseCase.CreateTickerUseCaseI(tickerPgRepo, tickerRedisRepo, cfg, appLogger)

		go func() {
			for {
				// Crawl tickers from vnd and save into database
				savedTickers, err := tickerUC.CrawlAllStockTicker(context.Background())
				if err != nil {
					appLogger.Fatal(err)
				}

				appLogger.Infof("Save %v ticker", len(savedTickers))

				time.Sleep(10 * time.Minute)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		sig := <-quit

		appLogger.Infof("Shutting down server... Reason: %s", sig)
	},
}

func init() {
	RootCmd.AddCommand(crawlSymbolCmd)
}

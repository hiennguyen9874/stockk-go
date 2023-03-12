package cmd

import (
	"context"
	"sync"
	"time"

	"github.com/hiennguyen9874/stockk-go/config"
	tickerRepository "github.com/hiennguyen9874/stockk-go/internal/tickers/repository"
	tickerUseCase "github.com/hiennguyen9874/stockk-go/internal/tickers/usecase"
	"github.com/hiennguyen9874/stockk-go/pkg/db/postgres"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/spf13/cobra"
)

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "crawl",
	Long:  "crawl",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetCfg()

		appLogger := logger.NewApiLogger(cfg)
		appLogger.InitLogger()
		appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.LoggerLevel, cfg.Server.Mode)

		psqlDB, err := postgres.NewPsqlDB(cfg)
		if err != nil {
			appLogger.Fatalf("Postgresql init: %s", err)
		} else {
			appLogger.Infof("Postgres connected")
		}

		// Repository
		tickerPgRepo := tickerRepository.CreateTickerPgRepository(psqlDB)

		// UseCase
		tickerUC := tickerUseCase.CreateTickerUseCaseI(tickerPgRepo, cfg, appLogger)

		for {
			var wg sync.WaitGroup

			wg.Add(1)

			go func() {
				defer wg.Done()

				// Crawl tickers from vnd and save into database
				savedTickers, err := tickerUC.CrawlAllStockTicker(context.Background())
				if err != nil {
					appLogger.Fatal(err)
				}

				appLogger.Infof("Save %v ticker", len(savedTickers))
			}()

			wg.Wait()

			time.Sleep(1 * time.Minute)
		}
	},
}

func init() {
	RootCmd.AddCommand(crawlCmd)
}

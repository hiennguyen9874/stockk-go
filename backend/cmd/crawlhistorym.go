package cmd

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/config"
	barRepository "github.com/hiennguyen9874/stockk-go/internal/bars/repository"
	barUseCase "github.com/hiennguyen9874/stockk-go/internal/bars/usecase"
	tickerRepository "github.com/hiennguyen9874/stockk-go/internal/tickers/repository"
	"github.com/hiennguyen9874/stockk-go/pkg/db/influxdb"
	"github.com/hiennguyen9874/stockk-go/pkg/db/postgres"
	"github.com/hiennguyen9874/stockk-go/pkg/db/redis"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/spf13/cobra"
)

var crawlHistoryMCmd = &cobra.Command{
	Use:   "crawlhistorym",
	Short: "crawl history",
	Long:  "crawl history",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

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

		influxDB, err := influxdb.NewInfluxDB(cfg)
		if err != nil {
			appLogger.Fatalf("InfluxDB init: %s", err)
		} else {
			appLogger.Infof("InfluxDB connected")
		}

		redisClient := redis.NewRedis(cfg)

		// Repository
		tickerPgRepo := tickerRepository.CreateTickerPgRepository(psqlDB)
		barInfluxDBRepo := barRepository.CreateBarRepo(influxDB, cfg.InfluxDB.InfluxDBOrg)
		barRedisRepo := barRepository.CreateBarRedisRepository(redisClient)

		barUseCase := barUseCase.CreateBarUseCaseI(barInfluxDBRepo, barRedisRepo, tickerPgRepo, cfg, appLogger)

		status, err := influxDB.Ping(ctx)
		if err != nil {
			appLogger.Warn(err)
			return
		}
		if !status {
			appLogger.Warn("influxdb not connected")
			return
		}

		appLogger.Info("Start syncing....")
		err = barUseCase.SyncMAllSymbol(ctx, cfg.Crawler.CrawlerTickerDownloadBatchSize, cfg.Crawler.CrawlerTickerInsertBatchSize, cfg.Crawler.CrawlerBarInsertBatchSize)
		if err != nil {
			appLogger.Warn(err)
		}
		appLogger.Info("Done sync, sleep 30s!")
		return
	},
}

func init() {
	RootCmd.AddCommand(crawlHistoryMCmd)
}

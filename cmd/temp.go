package cmd

import (
	"context"
	"time"

	"github.com/hiennguyen9874/stockk-go/config"
	barRepository "github.com/hiennguyen9874/stockk-go/internal/bars/repository"
	barUseCase "github.com/hiennguyen9874/stockk-go/internal/bars/usecase"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/pkg/crawlers"
	"github.com/hiennguyen9874/stockk-go/pkg/db/influxdb"
	"github.com/hiennguyen9874/stockk-go/pkg/db/redis"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/spf13/cobra"
)

var tempCmd = &cobra.Command{
	Use:   "temp",
	Short: "temp",
	Long:  "temp",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		cfg := config.GetCfg()

		appLogger := logger.NewApiLogger(cfg)
		appLogger.InitLogger()
		appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.LoggerLevel, cfg.Server.Mode)

		// psqlDB, err := postgres.NewPsqlDB(cfg)
		// if err != nil {
		// 	appLogger.Fatalf("Postgresql init: %s", err)
		// } else {
		// 	appLogger.Infof("Postgres connected")
		// }

		influxDB, err := influxdb.NewInfluxDB(cfg)
		if err != nil {
			appLogger.Fatalf("InfluxDB init: %s", err)
		} else {
			appLogger.Infof("InfluxDB connected")
		}

		redisClient := redis.NewRedis(cfg)

		status, err := influxDB.Ping(ctx)
		if err != nil {
			appLogger.Fatal(err)
		}

		appLogger.Info(status)

		// Repository
		// tickerPgRepo := tickerRepository.CreateTickerPgRepository(psqlDB)
		barInfluxDBRepo := barRepository.CreateBarRepo(influxDB, "history")
		barRedisRepo := barRepository.CreateBarRedisRepository(redisClient)

		barUseCase := barUseCase.CreateBarUseCaseI(barInfluxDBRepo, barRedisRepo, cfg)

		crawler := crawlers.NewCrawler(cfg)

		bucket := "long"

		from := time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC)
		to := time.Now().UTC()

		crawlerBars, err := crawler.VNDCrawlStockHistory("VCI", crawlers.RD, from.Unix(), to.Unix())
		if err != nil {
			appLogger.Fatal(err)
		}
		appLogger.Infof("Len bars: %v", len(crawlerBars))

		bars := make([]*models.Bar, len(crawlerBars))
		for i, crawlerBar := range crawlerBars {
			bars[i] = &models.Bar{
				Symbol:   "VCI",
				Exchange: "HSX",
				Time:     crawlerBar.Time,
				Open:     crawlerBar.Open,
				High:     crawlerBar.High,
				Low:      crawlerBar.Low,
				Close:    crawlerBar.Close,
				Volume:   crawlerBar.Volume,
			}
		}

		// err = barUseCase.Inserts(ctx, bucket, bars)
		// // err = barUseCase.Insert(ctx, bucket, bars[0])
		// if err != nil {
		// 	appLogger.Fatal(err)
		// }

		savedBars, err := barUseCase.GetByToLimit(ctx, bucket, "VCI", "HSX", to, 100)
		if err != nil {
			appLogger.Fatal(err)
		}

		appLogger.Info(savedBars)

		appLogger.Info("Done!")
	},
}

func init() {
	RootCmd.AddCommand(tempCmd)
}

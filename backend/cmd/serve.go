package cmd

import (
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/distributor"
	"github.com/hiennguyen9874/stockk-go/internal/server"
	"github.com/hiennguyen9874/stockk-go/pkg/db/influxdb"
	"github.com/hiennguyen9874/stockk-go/pkg/db/postgres"
	"github.com/hiennguyen9874/stockk-go/pkg/db/redis"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/hiennguyen9874/stockk-go/pkg/sentry"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Long:  `Starts a http server and serves the configured api`,
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

		if cfg.Server.MigrateOnStart {
			err = Migrate(psqlDB)

			if err != nil {
				appLogger.Info("Can not migrate data")
			} else {
				appLogger.Info("Data migrated")
			}
		}

		redisClient := redis.NewRedis(cfg)

		influxDB, err := influxdb.NewInfluxDB(cfg)
		if err != nil {
			appLogger.Fatalf("InfluxDB init: %s", err)
		} else {
			appLogger.Infof("InfluxDB connected")
		}

		taskRedisClient := distributor.NewRedisClient(cfg)

		server, err := server.NewServer(cfg, psqlDB, redisClient, taskRedisClient, influxDB, appLogger)
		if err != nil {
			appLogger.Fatal(err)
		}
		server.Start()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

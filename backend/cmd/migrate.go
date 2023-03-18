package cmd

import (
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/pkg/db/postgres"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate data",
	Long:  "Migrate data",
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

		err = Migrate(psqlDB)

		if err != nil {
			appLogger.Info("Can not migrate data")
		} else {
			appLogger.Info("Data migrated")
		}
	},
}

func Migrate(db *gorm.DB) error {
	var migrationModels = []interface{}{&models.User{}, &models.Ticker{}, &models.Chart{}, &models.StudyTemplate{}, &models.DrawingTemplate{}}

	err := db.AutoMigrate(migrationModels...)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}

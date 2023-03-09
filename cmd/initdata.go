package cmd

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/config"
	userRepository "github.com/hiennguyen9874/stockk-go/internal/users/repository"
	userUseCase "github.com/hiennguyen9874/stockk-go/internal/users/usecase"
	"github.com/hiennguyen9874/stockk-go/pkg/db/postgres"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/spf13/cobra"
)

var initDataCmd = &cobra.Command{
	Use:   "initdata",
	Short: "Init data",
	Long:  "Init data",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetCfg()

		appLogger := logger.NewApiLogger(cfg)
		appLogger.InitLogger()
		appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

		psqlDB, err := postgres.NewPsqlDB(cfg)
		if err != nil {
			appLogger.Fatalf("Postgresql init: %s", err)
		} else {
			appLogger.Infof("Postgres connected")
		}

		// Repository
		userRepo := userRepository.CreateUserRepository(psqlDB)

		// UseCase
		userUC := userUseCase.CreateUserUseCaseI(userRepo, cfg, appLogger)

		// Create super user if not exists
		isCreated, _ := userUC.CreateSuperUserIfNotExist(context.Background())

		if !isCreated {
			appLogger.Info("Super user is exists, skip create")
		} else {
			appLogger.Info("Created super user")
		}
	},
}

func init() {
	RootCmd.AddCommand(initDataCmd)
}

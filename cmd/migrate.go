package cmd

import (
	"log"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/pkg/db/postgres"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate data",
	Long:  "Migrate data",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetCfg()

		psqlDB, err := postgres.NewPsqlDB(cfg)
		if err != nil {
			log.Fatalf("Postgresql init: %s", err)
		} else {
			log.Println("Postgres connected")
		}

		Migrate(psqlDB)
	},
}

func Migrate(db *gorm.DB) {
	var migrationModels = []interface{}{&models.User{}}

	err := db.AutoMigrate(migrationModels...)
	if err != nil {
		return
	}
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}

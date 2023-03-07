package cmd

import (
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/models"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate data",
	Long:  "Migrate data",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := gorm.Open(postgres.Open(config.GetDNSConfig()), &gorm.Config{})

		if err == nil {
			Migrate(db)
		}
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

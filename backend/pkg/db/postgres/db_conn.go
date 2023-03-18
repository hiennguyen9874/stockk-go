package postgres

import (
	"fmt"

	"github.com/hiennguyen9874/stockk-go/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPsqlDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Postgres.PostgresqlHost, cfg.Postgres.PostgresqlUser, cfg.Postgres.PostgresqlPassword, cfg.Postgres.PostgresqlDbname, cfg.Postgres.PostgresqlPort, cfg.Postgres.PostgresqlSSLMode,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

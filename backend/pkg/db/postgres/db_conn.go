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
		cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Dbname, cfg.Postgres.Port, cfg.Postgres.SSLMode,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

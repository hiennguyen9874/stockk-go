package db

import (
	"github.com/hiennguyen9874/stockk-go/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetPostgres() (db *gorm.DB, err error) {
	return gorm.Open(postgres.Open(config.GetDNSConfig()), &gorm.Config{})
}

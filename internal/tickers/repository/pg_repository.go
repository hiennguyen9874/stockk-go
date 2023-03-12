package repository

import (
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/repository"
	"github.com/hiennguyen9874/stockk-go/internal/tickers"
	"gorm.io/gorm"
)

type TickerPgRepo struct {
	repository.PgRepo[models.Ticker]
}

func CreateTickerPgRepository(db *gorm.DB) tickers.TickerPgRepository {
	return &TickerPgRepo{
		PgRepo: repository.CreatePgRepo[models.Ticker](db),
	}
}

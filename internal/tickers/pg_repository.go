package tickers

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type TickerPgRepository interface {
	internal.PgRepository[models.Ticker]
	GetBySymbol(ctx context.Context, symbol string) (*models.Ticker, error)
	UpdateIsActive(ctx context.Context, exp *models.Ticker, isActive bool) (*models.Ticker, error)
	GetAllActive(ctx context.Context, isActive bool) ([]*models.Ticker, error)
}

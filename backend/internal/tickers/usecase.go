package tickers

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type TickerUseCaseI interface {
	internal.UseCaseI[models.Ticker]
	GetBySymbol(ctx context.Context, symbol string) (*models.Ticker, error)
	UpdateIsActiveBySymbol(ctx context.Context, symbol string, isActive bool) (*models.Ticker, error)
	CrawlAllStockTicker(ctx context.Context) ([]*models.Ticker, error)
	GetAllActive(ctx context.Context, isActive bool) ([]*models.Ticker, error)
	SearchBySymbol(ctx context.Context, symbol string, limit int, exchange string, isActive bool) ([]*models.Ticker, error)
}

package tickers

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type TickerUseCaseI interface {
	internal.UseCaseI[models.Ticker]
	CrawlAllStockTicker(ctx context.Context) ([]*models.Ticker, error)
}

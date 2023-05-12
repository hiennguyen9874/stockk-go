package stockssnapshot

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type StockSnapshotUseCaseI interface {
	CrawlAllStocksSnapshot(ctx context.Context) error
	GetStockSnapshotBySymbol(ctx context.Context, symbol string) (*models.StockSnapshot, error)
	UpdateStockSnapshotBySymbol(ctx context.Context, symbol string, values map[string]interface{}) error
}

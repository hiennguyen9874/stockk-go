package bars

import (
	"context"
	"time"

	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type BarUseCaseI interface {
	Insert(ctx context.Context, resolution string, exp *models.Bar, preventOverwriteOld bool) error
	Inserts(ctx context.Context, resolution string, exps []*models.Bar, barInsertBatchSize int, preventOverwriteOld bool) error
	GetByFromTo(ctx context.Context, resolution string, symbol string, from time.Time, to time.Time) ([]*models.Bar, error)
	GetByToLimit(ctx context.Context, resolution string, symbol string, to time.Time, limit int) ([]*models.Bar, error)
	SyncDSymbol(ctx context.Context, symbol string, barInsertBatchSize int) error
	SyncMSymbol(ctx context.Context, symbol string, barInsertBatchSize int) error
	SyncDAllSymbol(ctx context.Context, tickerDownloadBatchSize int, tickerInsertBatchSize int, barInsertBatchSize int) error
	SyncMAllSymbol(ctx context.Context, tickerDownloadBatchSize int, tickerInsertBatchSize int, barInsertBatchSize int) error
}

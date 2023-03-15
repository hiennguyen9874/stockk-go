package bars

import (
	"context"
	"time"

	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type BarUseCaseI interface {
	GenerateRedisLastTimestampKey(symbol string) string
	Insert(ctx context.Context, bucket string, exp *models.Bar) error
	Inserts(ctx context.Context, bucket string, exps []*models.Bar) error
	GetByFromTo(ctx context.Context, bucket string, exchange, symbol string, from time.Time, to time.Time) ([]*models.Bar, error)
	GetByToLimit(ctx context.Context, bucket string, symbol string, exchange string, to time.Time, limit int) ([]*models.Bar, error)
}

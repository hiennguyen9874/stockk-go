package bars

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type BarUseCaseI interface {
	GenerateRedisLastTimestampKey(symbol string) string
	Insert(ctx context.Context, bucket string, exp *models.Bar) error
	Inserts(ctx context.Context, bucket string, exps []*models.Bar) error
}

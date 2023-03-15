package bars

import (
	"context"
	"time"

	"github.com/hiennguyen9874/stockk-go/internal/models"
	influxdb2Write "github.com/influxdata/influxdb-client-go/v2/api/write"
)

type BarInfluxDBRepository interface {
	ToPoint(ctx context.Context, exp *models.Bar) *influxdb2Write.Point
	Insert(ctx context.Context, bucket string, exp *models.Bar) error
	Inserts(ctx context.Context, bucket string, exps []*models.Bar) error
	GetByFromTo(ctx context.Context, bucket string, symbol, exchange string, from time.Time, to time.Time) ([]*models.Bar, error)
	GetByToLimit(ctx context.Context, bucket string, symbol, exchange string, to time.Time, limit int, lastTime time.Time) ([]*models.Bar, error)
}

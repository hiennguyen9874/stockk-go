package bars

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal/models"
	influxdb2Write "github.com/influxdata/influxdb-client-go/v2/api/write"
)

type BarInfluxDBRepository interface {
	ToPoint(ctx context.Context, exp *models.Bar) *influxdb2Write.Point
	Create(ctx context.Context, bucket string, exp *models.Bar) error
}

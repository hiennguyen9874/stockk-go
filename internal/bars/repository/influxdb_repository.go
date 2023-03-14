package bars

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal/bars"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxdb2Write "github.com/influxdata/influxdb-client-go/v2/api/write"
)

type InfluxDBRepo struct {
	influxDBClient influxdb2.Client
	org            string
}

func CreateRedisRepo[M any](influxDBClient influxdb2.Client) bars.BarInfluxDBRepository {
	return &InfluxDBRepo{influxDBClient: influxDBClient}
}

func (r *InfluxDBRepo) ToPoint(ctx context.Context, exp *models.Bar) *influxdb2Write.Point {
	return influxdb2.NewPointWithMeasurement(exp.Exchange).
		AddTag("symbol", exp.Symbol).
		AddField("open", exp.Open).
		AddField("high", exp.High).
		AddField("low", exp.Low).
		AddField("close", exp.Close).
		AddField("volume", exp.Volume).
		SetTime(exp.Time)
}

func (r *InfluxDBRepo) Create(ctx context.Context, bucket string, exp *models.Bar) error {
	writeAPI := r.influxDBClient.WriteAPI(r.org, bucket)

	writeAPI.WritePoint(r.ToPoint(ctx, exp))

	return nil
}

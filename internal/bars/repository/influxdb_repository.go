package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/hiennguyen9874/stockk-go/internal/bars"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxdb2Write "github.com/influxdata/influxdb-client-go/v2/api/write"
)

type InfluxDBRepo struct {
	influxDBClient influxdb2.Client
	org            string
}

func CreateBarRepo(influxDBClient influxdb2.Client, org string) bars.BarInfluxDBRepository {
	return &InfluxDBRepo{influxDBClient: influxDBClient, org: org}
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

func (r *InfluxDBRepo) Insert(ctx context.Context, bucket string, exp *models.Bar) error {
	writeAPI := r.influxDBClient.WriteAPIBlocking(r.org, bucket)

	return writeAPI.WritePoint(ctx, r.ToPoint(ctx, exp))
}

func (r *InfluxDBRepo) Inserts(ctx context.Context, bucket string, exps []*models.Bar) error {
	writeAPI := r.influxDBClient.WriteAPIBlocking(r.org, bucket)

	var wg sync.WaitGroup

	for _, exp := range exps {
		wg.Add(1)

		go func(exp *models.Bar) {
			defer wg.Done()
			err := writeAPI.WritePoint(ctx, r.ToPoint(ctx, exp))
			if err != nil {
				fmt.Print(err)
			}
		}(exp)
	}
	wg.Wait()
	return nil
}

// func (r *InfluxDBRepo) GetLastBySymbol(ctx context.Context, bucket string, symbol string, exchange string) (*models.Bar, error) {
// 	queryAPI := r.influxDBClient.QueryAPI(r.org)

// 	// Query
// 	query := fmt.Sprintf(`from(bucket:"%v")
// 	|> range(start: -10d)
// 	|> filter(fn: (r) => r._measurement == "%v")
// 	|> filter(fn: (r) => r.symbol == "%v")
// 	|> last()`, bucket, exchange, symbol)

// 	fmt.Println(query)

// 	// Get result
// 	result, err := queryAPI.Query(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// fmt.Println(result)

// 	for result.Next() {
// 		// Access data
// 		fmt.Printf("value: %v\n", result.Record().Value())
// 	}
// 	return nil, errors.New("not implemented error")
// }

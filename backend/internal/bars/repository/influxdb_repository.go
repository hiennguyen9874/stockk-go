package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/hiennguyen9874/stockk-go/internal/bars"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxdb2API "github.com/influxdata/influxdb-client-go/v2/api"
	influxdb2Write "github.com/influxdata/influxdb-client-go/v2/api/write"
)

type BarInfluxDBRepo struct {
	influxDBClient influxdb2.Client
	org            string
}

func CreateBarRepo(influxDBClient influxdb2.Client, org string) bars.BarInfluxDBRepository {
	return &BarInfluxDBRepo{influxDBClient: influxDBClient, org: org}
}

func (r *BarInfluxDBRepo) ToPoint(ctx context.Context, exp *models.Bar) *influxdb2Write.Point {
	return influxdb2.NewPointWithMeasurement(exp.Exchange).
		AddTag("symbol", exp.Symbol).
		AddField("open", exp.Open).
		AddField("high", exp.High).
		AddField("low", exp.Low).
		AddField("close", exp.Close).
		AddField("volume", exp.Volume).
		SetTime(exp.Time)
}

func (r *BarInfluxDBRepo) Insert(ctx context.Context, bucket string, exp *models.Bar) error {
	writeAPI := r.influxDBClient.WriteAPIBlocking(r.org, bucket)

	return writeAPI.WritePoint(ctx, r.ToPoint(ctx, exp))
}

func (r *BarInfluxDBRepo) Inserts(ctx context.Context, bucket string, exps []*models.Bar, batchSize int) error {
	writeAPI := r.influxDBClient.WriteAPIBlocking(r.org, bucket)

	// var wg sync.WaitGroup

	// for _, exp := range exps {
	// 	wg.Add(1)

	// 	go func(exp *models.Bar) {
	// 		defer wg.Done()
	// 		err := writeAPI.WritePoint(ctx, r.ToPoint(ctx, exp))
	// 		if err != nil {
	// 			fmt.Print(err)
	// 		}
	// 	}(exp)
	// }
	// wg.Wait()

	// for _, exp := range exps {
	// 	err := writeAPI.WritePoint(ctx, r.ToPoint(ctx, exp))
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// Ticker queue
	expsCh := make(chan *models.Bar, batchSize)
	doneCh := make(chan bool)
	errCh := make(chan error)

	// Ticker producer
	go func() {
		var wg sync.WaitGroup
		for _, bar := range exps {
			wg.Add(1)
			go func(expsCh chan<- *models.Bar, ticker *models.Bar) {
				defer wg.Done()
				expsCh <- ticker
			}(expsCh, bar)
		}
		wg.Wait()
		close(expsCh)
	}()

	// Ticker consumer
	go func(expsCh <-chan *models.Bar, doneCh chan<- bool, errCh chan<- error) {
		var wg sync.WaitGroup
		for exp := range expsCh {
			wg.Add(1)
			go func(exp *models.Bar) {
				defer wg.Done()

				err := writeAPI.WritePoint(ctx, r.ToPoint(ctx, exp))
				if err != nil {
					errCh <- err
				}
			}(exp)
		}
		wg.Wait()
		doneCh <- true
	}(expsCh, doneCh, errCh)

	select {
	case err := <-errCh:
		return err
	case <-doneCh:
		return nil
	}
}

func (r *BarInfluxDBRepo) ParseResultFromInfluxDB(result *influxdb2API.QueryTableResult) ([]*models.Bar, error) {
	var bars []*models.Bar

	type Record struct {
		Symbol   string
		Exchange string
		Open     float64
		High     float64
		Low      float64
		Close    float64
		Volume   int64
	}

	records := make(map[time.Time]Record)
	for result.Next() {
		val, ok := records[result.Record().Time()]
		if !ok {
			val = Record{}
			val.Symbol = result.Record().ValueByKey("symbol").(string)
			val.Exchange = result.Record().Measurement()
		}

		switch field := result.Record().Field(); field {
		case "open":
			val.Open = result.Record().Value().(float64)
		case "high":
			val.High = result.Record().Value().(float64)
		case "low":
			val.Low = result.Record().Value().(float64)
		case "close":
			val.Close = result.Record().Value().(float64)
		case "volume":
			val.Volume = result.Record().Value().(int64)
		}
		records[result.Record().Time()] = val
	}

	for recordTime, record := range records {
		bars = append(bars, &models.Bar{
			Symbol:   record.Symbol,
			Exchange: record.Exchange,
			Time:     recordTime,
			Open:     record.Open,
			High:     record.High,
			Low:      record.Low,
			Close:    record.Close,
			Volume:   record.Volume,
		})
	}

	return bars, nil
}

func (r *BarInfluxDBRepo) GetByFromTo(ctx context.Context, bucket string, symbol, exchange string, from time.Time, to time.Time) ([]*models.Bar, error) {
	queryAPI := r.influxDBClient.QueryAPI(r.org)

	// Query
	query := fmt.Sprintf(`from(bucket:"%v")
		|> range(start: %v, stop: %v)
		|> filter(fn: (r) => r._measurement == "%v")
		|> filter(fn: (r) => r.symbol == "%v")
		|> sort(columns: ["_time"], desc: false)`, bucket, from.Unix(), to.Unix(), exchange, symbol)

	// Get result
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	return r.ParseResultFromInfluxDB(result)
}

func (r *BarInfluxDBRepo) GetByToLimit(ctx context.Context, bucket string, symbol, exchange string, to time.Time, limit int, lastTime time.Time) ([]*models.Bar, error) {
	queryAPI := r.influxDBClient.QueryAPI(r.org)

	startTime, err := r.estimateStart(lastTime, limit, "D")
	if err != nil {
		return nil, err
	}

	// Query
	query := fmt.Sprintf(`from(bucket:"%v")
		|> range(start: %v, stop: %v)
		|> filter(fn: (r) => r._measurement == "%v")
		|> filter(fn: (r) => r.symbol == "%v")
		|> sort(columns: ["_time"], desc: true)
    	|> limit(n: %v)
		|> sort(columns: ["_time"], desc: false)`, bucket, startTime.Unix(), to.Unix(), exchange, symbol, limit)

	// Get result
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	return r.ParseResultFromInfluxDB(result)
}

func (r *BarInfluxDBRepo) estimateStart(to time.Time, limit int, resolution string) (time.Time, error) {
	switch strings.ToLower(resolution) {
	case "d":
		return to.AddDate(0, 0, -int(math.Ceil(float64(limit)*7/4))), nil
	default:
		// TODO: Implement for 1m
		return time.Time{}, errors.New("not implemented error")
	}
}

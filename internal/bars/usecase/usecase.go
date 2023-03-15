package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/bars"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/pkg/crawlers"
)

type barUseCase struct {
	influxDBRepo bars.BarInfluxDBRepository
	redisRepo    bars.BarRedisRepository
	crawler      crawlers.Crawler
}

func CreateBarUseCaseI(
	influxDBRepo bars.BarInfluxDBRepository,
	redisRepo bars.BarRedisRepository,
	cfg *config.Config,
) bars.BarUseCaseI {
	return &barUseCase{
		influxDBRepo: influxDBRepo,
		redisRepo:    redisRepo,
		crawler:      crawlers.NewCrawler(cfg),
	}
}

func (u *barUseCase) Insert(ctx context.Context, bucket string, exp *models.Bar) error {
	// Get last timestamp from redis
	timestamp, err := u.redisRepo.GetInt64(ctx, u.GenerateRedisLastTimestampKey(exp.Symbol))
	if err != nil {
		return err
	}

	// Check timestamp: Only insert new timestamp into influxdb
	if timestamp != nil && exp.Time.Unix() < *timestamp {
		return errors.New("can not insert old timestamp")
	}

	// Save timestamp in redis to last timestamp
	if timestamp == nil || (timestamp != nil && exp.Time.Unix() > *timestamp) {
		err = u.redisRepo.CreateInt64(ctx, u.GenerateRedisLastTimestampKey(exp.Symbol), exp.Time.Unix(), -1)

		if err != nil {
			return err
		}
	}

	// Insert into influxdb
	return u.influxDBRepo.Insert(ctx, bucket, exp)
}

func (u *barUseCase) Inserts(ctx context.Context, bucket string, exps []*models.Bar) error {
	savedLastTimestamp := make(map[string]int64)
	for _, exp := range exps {
		// Get last timestamp from redis
		timestamp, err := u.redisRepo.GetInt64(ctx, u.GenerateRedisLastTimestampKey(exp.Symbol))
		if err != nil {
			return err
		}

		if timestamp != nil {
			savedLastTimestamp[exp.Symbol] = *timestamp
		} else {
			savedLastTimestamp[exp.Symbol] = -1
		}
	}

	for _, exp := range exps {
		// Check timestamp: Only insert new timestamp into influxdb
		if exp.Time.Unix() < savedLastTimestamp[exp.Symbol] {
			return errors.New("can not insert old timestamp")
		}
	}

	// Get last timestamp
	lastTimestamp := make(map[string]int64)
	for _, exp := range exps {
		if exp.Time.Unix() > lastTimestamp[exp.Symbol] {
			lastTimestamp[exp.Symbol] = exp.Time.Unix()
		}
	}

	// Save last timestamp into redis
	for symbol, timestamp := range lastTimestamp {
		err := u.redisRepo.CreateInt64(ctx, u.GenerateRedisLastTimestampKey(symbol), timestamp, -1)
		if err != nil {
			return err
		}
	}

	// Insert into influxdb
	return u.influxDBRepo.Inserts(ctx, bucket, exps)
}

func (u *barUseCase) GetByFromTo(ctx context.Context, bucket string, exchange, symbol string, from time.Time, to time.Time) ([]*models.Bar, error) {
	return u.influxDBRepo.GetByFromTo(ctx, bucket, symbol, exchange, from, to)
}

func (u *barUseCase) GetByToLimit(ctx context.Context, bucket string, symbol string, exchange string, to time.Time, limit int) ([]*models.Bar, error) {
	timestamp, err := u.redisRepo.GetInt64(ctx, u.GenerateRedisLastTimestampKey(symbol))
	if err != nil {
		return nil, err
	}

	if timestamp == nil {
		return nil, errors.New("last timestamp not saved")
	}

	return u.influxDBRepo.GetByToLimit(ctx, bucket, symbol, exchange, to, limit, time.Unix(*timestamp, 0))
}

func (u *barUseCase) GenerateRedisLastTimestampKey(symbol string) string {
	return fmt.Sprintf("%s:%s", "LastTimeStamp", symbol)
}

package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/bars"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/tickers"
	"github.com/hiennguyen9874/stockk-go/pkg/crawlers"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

type barUseCase struct {
	barInfluxDBRepo bars.BarInfluxDBRepository
	barRedisRepo    bars.BarRedisRepository
	tickerPgRepo    tickers.TickerPgRepository
	crawler         crawlers.Crawler
	logger          logger.Logger
}

func CreateBarUseCaseI(
	barInfluxDBRepo bars.BarInfluxDBRepository,
	barRedisRepo bars.BarRedisRepository,
	tickerPgRepo tickers.TickerPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) bars.BarUseCaseI {
	return &barUseCase{
		barInfluxDBRepo: barInfluxDBRepo,
		barRedisRepo:    barRedisRepo,
		tickerPgRepo:    tickerPgRepo,
		crawler:         crawlers.NewCrawler(cfg, logger),
		logger:          logger,
	}
}

func (u *barUseCase) convertResolutionToBucket(ctx context.Context, resolution string) (string, error) {
	switch resolution {
	case "1":
		return "Resolution1", nil
	case "5":
		return "Resolution5", nil
	case "15":
		return "Resolution15", nil
	case "30":
		return "Resolution30", nil
	case "60":
		return "Resolution60", nil
	case "D":
		return "ResolutionD", nil
	default:
		// TODO: Use httpErrors
		return "", errors.New("can not convert resolution to bucket")
	}
}

func (u *barUseCase) convertBucketToResolution(ctx context.Context, bucket string) (string, error) {
	switch bucket {
	case "Resolution1":
		return "1", nil
	case "Resolution5":
		return "5", nil
	case "Resolution15":
		return "15", nil
	case "Resolution30":
		return "30", nil
	case "Resolution60":
		return "60", nil
	case "ResolutionD":
		return "D", nil
	default:
		// TODO: Use httpErrors
		return "", errors.New("not support resolution")
	}
}

func (u *barUseCase) generateRedisLastTimestampKey(symbol string, resolution string) string {
	return fmt.Sprintf("%s:%s:%s", "LastTimeStamp", resolution, symbol)
}

func (u *barUseCase) convertResolutionToCrawlerResolution(resolution string) (crawlers.Resolution, error) {
	switch resolution {
	case "1":
		return crawlers.R1, nil
	case "5":
		return crawlers.R5, nil
	case "15":
		return crawlers.R15, nil
	case "30":
		return crawlers.R30, nil
	case "60":
		return crawlers.R60, nil
	case "D":
		return crawlers.RD, nil
	default:
		// TODO: Use httpErrors
		return crawlers.RD, fmt.Errorf("not support resolution: %v", resolution)
	}
}

func (u *barUseCase) convertResolutionToTimeDuration(resolution string) (time.Duration, error) {
	switch resolution {
	case "1":
		return time.Duration(time.Minute), nil
	case "5":
		return time.Duration(time.Minute * 5), nil
	case "15":
		return time.Duration(time.Minute * 15), nil
	case "30":
		return time.Duration(time.Minute * 30), nil
	case "60":
		return time.Duration(time.Hour), nil
	case "D":
		return time.Duration(time.Hour * 24), nil
	default:
		// TODO: Use httpErrors
		return time.Duration(time.Hour * 24), fmt.Errorf("not support resolution: %v", resolution)
	}
}

func (u *barUseCase) Insert(ctx context.Context, resolution string, exp *models.Bar, preventOverwriteOld bool) error {
	// Get last timestamp from redis
	timestamp, err := u.barRedisRepo.GetInt64(ctx, u.generateRedisLastTimestampKey(exp.Symbol, resolution))
	if err != nil {
		return err
	}

	if preventOverwriteOld {
		// Check timestamp: Only insert new timestamp into influxdb
		if timestamp != nil && exp.Time.Unix() < *timestamp {
			return errors.New("can not insert old timestamp")
		}
	}

	// Save timestamp in redis to last timestamp
	if timestamp == nil || (timestamp != nil && exp.Time.Unix() > *timestamp) {
		err = u.barRedisRepo.CreateInt64(ctx, u.generateRedisLastTimestampKey(exp.Symbol, resolution), exp.Time.Unix(), -1)

		if err != nil {
			return err
		}
	}

	bucket, err := u.convertResolutionToBucket(ctx, resolution)
	if err != nil {
		return err
	}

	// Insert into influxdb
	return u.barInfluxDBRepo.Insert(ctx, bucket, exp)
}

func (u *barUseCase) Inserts(ctx context.Context, resolution string, exps []*models.Bar, barInsertBatchSize int, preventOverwriteOld bool) error {
	savedLastTimestamp := make(map[string]int64)
	for _, exp := range exps {
		// Get last timestamp from redis
		timestamp, err := u.barRedisRepo.GetInt64(ctx, u.generateRedisLastTimestampKey(exp.Symbol, resolution))
		if err != nil {
			return err
		}

		if timestamp != nil {
			savedLastTimestamp[exp.Symbol] = *timestamp
		} else {
			savedLastTimestamp[exp.Symbol] = -1
		}
	}

	if preventOverwriteOld {
		for _, exp := range exps {
			// Check timestamp: Only insert new timestamp into influxdb
			if exp.Time.Unix() < savedLastTimestamp[exp.Symbol] {
				return errors.New("can not insert old timestamp")
			}
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
		err := u.barRedisRepo.CreateInt64(ctx, u.generateRedisLastTimestampKey(symbol, resolution), timestamp, -1)
		if err != nil {
			return err
		}
	}

	bucket, err := u.convertResolutionToBucket(ctx, resolution)
	if err != nil {
		return err
	}

	// Insert into influxdb
	return u.barInfluxDBRepo.Inserts(ctx, bucket, exps, barInsertBatchSize)
}

func (u *barUseCase) GetByFromTo(ctx context.Context, resolution string, symbol string, from time.Time, to time.Time) ([]*models.Bar, error) {
	ticker, err := u.tickerPgRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	bucket, err := u.convertResolutionToBucket(ctx, resolution)
	if err != nil {
		return nil, err
	}

	return u.barInfluxDBRepo.GetByFromTo(ctx, bucket, ticker.Symbol, ticker.Exchange, from, to)
}

func (u *barUseCase) GetByToLimit(ctx context.Context, resolution string, symbol string, to time.Time, limit int) ([]*models.Bar, error) {
	ticker, err := u.tickerPgRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	timestamp, err := u.barRedisRepo.GetInt64(ctx, u.generateRedisLastTimestampKey(ticker.Symbol, resolution))
	if err != nil {
		return nil, err
	}

	if timestamp == nil {
		return nil, errors.New("last timestamp not saved")
	}

	bucket, err := u.convertResolutionToBucket(ctx, resolution)
	if err != nil {
		return nil, err
	}

	return u.barInfluxDBRepo.GetByToLimit(ctx, bucket, ticker.Symbol, ticker.Exchange, to, limit, time.Unix(*timestamp, 0))
}

func (u *barUseCase) CrawlSymbol(ctx context.Context, symbol string, resolution string, from time.Time, to time.Time, barInsertBatchSize int) error {
	ticker, err := u.tickerPgRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return err
	}

	crawlerResolution, err := u.convertResolutionToCrawlerResolution(resolution)
	if err != nil {
		return err
	}

	crawlerBars, err := u.crawler.CrawlStockHistory(ctx, ticker.Symbol, crawlerResolution, from.Unix(), to.Unix())
	if err != nil {
		return err
	}

	bars := make([]*models.Bar, len(crawlerBars))
	for i, crawlerBar := range crawlerBars {
		bars[i] = &models.Bar{
			Symbol:   ticker.Symbol,
			Exchange: ticker.Exchange,
			Time:     crawlerBar.Time,
			Open:     crawlerBar.Open,
			High:     crawlerBar.High,
			Low:      crawlerBar.Low,
			Close:    crawlerBar.Close,
			Volume:   crawlerBar.Volume,
		}
	}

	return u.Inserts(ctx, resolution, bars, barInsertBatchSize, false)
}

func (u *barUseCase) SyncSymbol(ctx context.Context, symbol string, resolution string, barInsertBatchSize int) error {
	ticker, err := u.tickerPgRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return err
	}

	var fromTime time.Time

	// Get last timestamp from redis
	timestamp, err := u.barRedisRepo.GetInt64(ctx, u.generateRedisLastTimestampKey(ticker.Symbol, resolution))
	if err != nil {
		return err
	}

	if timestamp != nil {
		fromTime = time.Unix(*timestamp, 0)
	} else {
		fromTime = time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC)
	}

	toTime := time.Now().UTC()

	return u.CrawlSymbol(ctx, symbol, resolution, fromTime, toTime, barInsertBatchSize)
}

func (u *barUseCase) SyncAllSymbol(ctx context.Context, resolution string, tickerDownloadBatchSize int, tickerInsertBatchSize int, barInsertBatchSize int) error {
	activeTickers, err := u.tickerPgRepo.GetAllActive(ctx, true)
	if err != nil {
		return err
	}

	if len(activeTickers) == 0 {
		return errors.New("not found any ticker active")
	}

	crawlerResolution, err := u.convertResolutionToCrawlerResolution(resolution)
	if err != nil {
		return err
	}

	// TODO: Add context.Done() into goroutine

	// Ticker queue
	tickersCh := make(chan *models.Ticker, tickerDownloadBatchSize)

	go func() {
		var queueWg sync.WaitGroup
		for _, activeTicker := range activeTickers {
			queueWg.Add(1)
			go func(tickersCh chan<- *models.Ticker, ticker *models.Ticker) {
				defer queueWg.Done()
				tickersCh <- ticker
			}(tickersCh, activeTicker)
		}
		queueWg.Wait()
		close(tickersCh)
	}()

	// Download
	barsCh := make(chan []*models.Bar, tickerInsertBatchSize)
	errSendCh := make(chan error)

	go func(barsCh chan<- []*models.Bar, errCh chan<- error, resolution string) {
		var sendWg sync.WaitGroup
		for ticker := range tickersCh {
			sendWg.Add(1)
			go func(barsCh chan<- []*models.Bar, errCh chan<- error, ticker *models.Ticker, resolution string) {
				defer sendWg.Done()

				var fromTime time.Time

				// Get last timestamp from redis
				timestamp, err := u.barRedisRepo.GetInt64(ctx, u.generateRedisLastTimestampKey(ticker.Symbol, resolution))
				if err != nil {
					errCh <- err
				}

				if timestamp != nil {
					fromTime = time.Unix(*timestamp, 0)
				} else {
					fromTime = time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC)
				}

				toTime := time.Now().UTC()

				crawlerBars, err := u.crawler.CrawlStockHistory(ctx, ticker.Symbol, crawlerResolution, fromTime.Unix(), toTime.Unix())
				if err != nil {
					errCh <- err
				}

				bars := make([]*models.Bar, len(crawlerBars))
				for i, crawlerBar := range crawlerBars {
					bars[i] = &models.Bar{
						Symbol:   ticker.Symbol,
						Exchange: ticker.Exchange,
						Time:     crawlerBar.Time,
						Open:     crawlerBar.Open,
						High:     crawlerBar.High,
						Low:      crawlerBar.Low,
						Close:    crawlerBar.Close,
						Volume:   crawlerBar.Volume,
					}
				}

				barsCh <- bars
			}(barsCh, errSendCh, ticker, resolution)
		}
		sendWg.Wait()
		close(barsCh)
	}(barsCh, errSendCh, resolution)

	// Save
	doneCh := make(chan bool)
	errReceiveCh := make(chan error)

	go func(barsCh <-chan []*models.Bar, doneCh chan<- bool, errCh chan<- error) {
		var receiveWg sync.WaitGroup
		for bars := range barsCh {
			receiveWg.Add(1)
			go func(bars []*models.Bar) {
				defer receiveWg.Done()

				err := u.Inserts(ctx, resolution, bars, barInsertBatchSize, false)
				if err != nil {
					errCh <- err
				}
			}(bars)
		}
		receiveWg.Wait()
		doneCh <- true
	}(barsCh, doneCh, errReceiveCh)

	select {
	case err := <-errSendCh:
		return err
	case err := <-errReceiveCh:
		return err
	case <-doneCh:
		return nil
	}
}

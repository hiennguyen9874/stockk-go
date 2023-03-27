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

// func (u *barUseCase) convertBucketToResolution(ctx context.Context, bucket string) (string, error) {
// 	switch bucket {
// 	case "Resolution1":
// 		return "1", nil
// 	case "Resolution5":
// 		return "5", nil
// 	case "Resolution15":
// 		return "15", nil
// 	case "Resolution30":
// 		return "30", nil
// 	case "Resolution60":
// 		return "60", nil
// 	case "ResolutionD":
// 		return "D", nil
// 	default:
// 		// TODO: Use httpErrors
// 		return "", errors.New("not support resolution")
// 	}
// }

func (u *barUseCase) genRedisLastTimestampKey(symbol string, resolution string) string {
	return fmt.Sprintf("%s:%s:%s", "LastTimeStamp", resolution, symbol)
}

func (u *barUseCase) genRedisLastBar(symbol string, resolution string) string {
	return fmt.Sprintf("%s:%s:%s", "LastBar", resolution, symbol)
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

// func (u *barUseCase) convertResolutionToTimeDuration(resolution string) (time.Duration, error) {
// 	switch resolution {
// 	case "1":
// 		return time.Duration(time.Minute), nil
// 	case "5":
// 		return time.Duration(time.Minute * 5), nil
// 	case "15":
// 		return time.Duration(time.Minute * 15), nil
// 	case "30":
// 		return time.Duration(time.Minute * 30), nil
// 	case "60":
// 		return time.Duration(time.Hour), nil
// 	case "D":
// 		return time.Duration(time.Hour * 24), nil
// 	default:
// 		// TODO: Use httpErrors
// 		return time.Duration(time.Hour * 24), fmt.Errorf("not support resolution: %v", resolution)
// 	}
// }

func (u *barUseCase) Insert(ctx context.Context, resolution string, exp *models.Bar, preventOverwriteOld bool) error {
	// Get last bar from redis
	lastBar, err := u.barRedisRepo.GetObj(ctx, u.genRedisLastBar(exp.Symbol, resolution))
	if err != nil {
		return err
	}

	if preventOverwriteOld {
		// Check bar: Only insert new bar into influxdb
		if lastBar != nil && exp.Time.Unix() < lastBar.Time.Unix() {
			return errors.New("can not insert old timestamp")
		}
	}

	// Save last bar into redis
	if lastBar == nil || (lastBar != nil && exp.Time.Unix() > lastBar.Time.Unix()) {
		err = u.barRedisRepo.CreateObj(ctx, u.genRedisLastBar(exp.Symbol, resolution), exp, -1)

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
	symbols := make(map[string]bool)
	for _, exp := range exps {
		symbols[exp.Symbol] = true
	}

	savedLastBar := make(map[string]*models.Bar)
	for symbol := range symbols {
		// Get last bar from redis
		lastBar, err := u.barRedisRepo.GetObj(ctx, u.genRedisLastBar(symbol, resolution))
		if err != nil {
			return err
		}

		if lastBar != nil {
			savedLastBar[symbol] = lastBar
		} else {
			savedLastBar[symbol] = nil
		}
	}

	if preventOverwriteOld {
		for _, exp := range exps {
			// Check bar: Only insert new bar into influxdb
			if savedLastBar[exp.Symbol] == nil || exp.Time.Unix() < savedLastBar[exp.Symbol].Time.Unix() {
				return errors.New("can not insert old bar")
			}
		}
	}

	// Get last bar
	lastBar := make(map[string]*models.Bar)
	for _, exp := range exps {
		if savedLastBar[exp.Symbol] == nil {
			lastBar[exp.Symbol] = exp
		}

		if lastBar[exp.Symbol] == nil {
			lastBar[exp.Symbol] = exp
		} else if exp.Time.Unix() > lastBar[exp.Symbol].Time.Unix() {
			lastBar[exp.Symbol] = exp
		}
	}

	// Save last bar into redis
	for symbol, bar := range lastBar {
		err := u.barRedisRepo.CreateObj(ctx, u.genRedisLastBar(symbol, resolution), bar, -1)
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

	lastBar, err := u.barRedisRepo.GetObj(ctx, u.genRedisLastBar(ticker.Symbol, resolution))
	if err != nil {
		return nil, err
	}

	if lastBar == nil {
		return nil, errors.New("last bar not saved")
	}

	bucket, err := u.convertResolutionToBucket(ctx, resolution)
	if err != nil {
		return nil, err
	}

	return u.barInfluxDBRepo.GetByToLimit(ctx, bucket, ticker.Symbol, ticker.Exchange, to, limit, lastBar.Time)
}

func (u *barUseCase) SyncDSymbol(ctx context.Context, symbol string, barInsertBatchSize int) error {
	resolution := "D"

	ticker, err := u.tickerPgRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return err
	}

	var fromTime time.Time

	// Get last bar from redis
	lastBar, err := u.barRedisRepo.GetObj(ctx, u.genRedisLastBar(ticker.Symbol, resolution))
	if err != nil {
		return err
	}

	if lastBar != nil {
		fromTime = lastBar.Time
	} else {
		fromTime = time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC)
	}

	toTime := time.Now().UTC()

	crawlerResolution, err := u.convertResolutionToCrawlerResolution(resolution)
	if err != nil {
		return err
	}

	crawlerBars, err := u.crawler.CrawlStockHistory(ctx, ticker.Symbol, crawlerResolution, fromTime.Unix(), toTime.Unix())
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

	err = u.Inserts(ctx, resolution, bars, barInsertBatchSize, false)
	if err != nil {
		return err
	}

	return err
}

func (u *barUseCase) syncDAllSymbol(ctx context.Context, tickerDownloadBatchSize int, tickerInsertBatchSize int, barInsertBatchSize int, activeTickers []*models.Ticker) error {
	resolution := "D"

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

				// Get last bar from redis
				lastBar, err := u.barRedisRepo.GetObj(ctx, u.genRedisLastBar(ticker.Symbol, resolution))
				if err != nil {
					errCh <- err
				}

				if lastBar != nil {
					fromTime = lastBar.Time
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

func (u *barUseCase) SyncDAllSymbol(ctx context.Context, tickerDownloadBatchSize int, tickerInsertBatchSize int, barInsertBatchSize int) error {
	activeTickers, err := u.tickerPgRepo.GetAllActive(ctx, true)
	if err != nil {
		return err
	}

	if len(activeTickers) == 0 {
		return errors.New("not found any ticker active")
	}

	return u.syncDAllSymbol(ctx, tickerDownloadBatchSize, tickerInsertBatchSize, barInsertBatchSize, activeTickers)
}

func (u *barUseCase) SyncMSymbol(ctx context.Context, symbol string, barInsertBatchSize int) error {
	resolution := "1"
	resolutionD := "D"

	ticker, err := u.tickerPgRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return err
	}

	var fromTime time.Time

	// Get last bar from redis
	lastBar, err := u.barRedisRepo.GetObj(ctx, u.genRedisLastBar(ticker.Symbol, resolution))
	if err != nil {
		return err
	}

	if lastBar != nil {
		fromTime = lastBar.Time
	} else {
		fromTime = time.Now().UTC().Add(-time.Duration(30) * time.Hour * 24)
	}

	toTime := time.Now().UTC()

	crawlerResolution, err := u.convertResolutionToCrawlerResolution(resolution)
	if err != nil {
		return err
	}

	crawlerBars, err := u.crawler.CrawlStockHistory(ctx, ticker.Symbol, crawlerResolution, fromTime.Unix(), toTime.Unix())
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

	err = u.Inserts(ctx, resolution, bars, barInsertBatchSize, false)
	if err != nil {
		return err
	}

	if lastBar == nil {
		return err
	}

	var lastMBar *models.Bar = bars[0]
	for _, bar := range bars {
		if bar.Time.Unix() > lastMBar.Time.Unix() {
			lastMBar = bar
		}
	}

	if lastMBar.Time.Format("01-02-2006") != lastBar.Time.Format("01-02-2006") {
		return err
	}

	// Get last bar from redis
	lastDBar, err := u.barRedisRepo.GetObj(ctx, u.genRedisLastBar(ticker.Symbol, resolutionD))
	if err != nil {
		return err
	}

	if lastMBar.Time.Format("01-02-2006") != lastDBar.Time.Format("01-02-2006") {
		return u.SyncDSymbol(ctx, symbol, barInsertBatchSize)
	}

	newLastDBar := &models.Bar{
		Symbol:   lastDBar.Symbol,
		Exchange: lastDBar.Exchange,
		Time:     lastDBar.Time,
		Open:     lastDBar.Open,
		High:     lastDBar.High,
		Low:      lastDBar.Low,
		Close:    lastDBar.Close,
		Volume:   lastDBar.Volume,
	}
	hasNew := false

	for _, bar := range bars {
		// TODO: Convert to symbol timezone before compare
		if bar.Time.Format("01-02-2006") == newLastDBar.Time.Format("01-02-2006") {
			if bar.High > newLastDBar.High {
				newLastDBar.High = bar.High
				hasNew = true
			}

			if bar.Low < newLastDBar.Low {
				newLastDBar.Low = bar.Low
				hasNew = true
			}

			if lastBar.Time.Unix() < bar.Time.Unix() {
				newLastDBar.Volume += bar.Volume
				hasNew = true
			}
		}
	}

	if hasNew {
		return u.Insert(ctx, resolutionD, newLastDBar, false)
	}

	return nil
}

func (u *barUseCase) SyncMAllSymbol(ctx context.Context, tickerDownloadBatchSize int, tickerInsertBatchSize int, barInsertBatchSize int) error {
	resolution := "1"
	resolutionD := "D"

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
	newLastDBarCh := make(chan *models.Bar, tickerInsertBatchSize)
	mustSyncDTicker := make(chan *models.Ticker)
	errSendCh := make(chan error)

	go func(barsCh chan<- []*models.Bar, newLastDBarCh chan<- *models.Bar, mustSyncDTicker chan<- *models.Ticker, errCh chan<- error) {
		var sendWg sync.WaitGroup
		for ticker := range tickersCh {
			sendWg.Add(1)
			go func(barsCh chan<- []*models.Bar, newLastDBarCh chan<- *models.Bar, mustSyncDTicker chan<- *models.Ticker, errCh chan<- error, ticker *models.Ticker) {
				defer sendWg.Done()

				var fromTime time.Time

				// Get last bar from redis
				lastBar, err := u.barRedisRepo.GetObj(ctx, u.genRedisLastBar(ticker.Symbol, resolution))
				if err != nil {
					errCh <- err
				}

				if lastBar != nil {
					fromTime = lastBar.Time.Add(-time.Duration(1) * time.Minute)
				} else {
					fromTime = time.Now().UTC().Add(-6 * 30 * 24 * time.Hour)
				}

				toTime := time.Now().UTC()

				crawlerBars, err := u.crawler.CrawlStockHistory(ctx, ticker.Symbol, crawlerResolution, fromTime.Unix(), toTime.Unix())
				if err != nil {
					errCh <- err
				}

				if len(crawlerBars) == 0 {
					return
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

				if lastBar == nil {
					return
				}

				var lastMBar *models.Bar = bars[0]
				for _, bar := range bars {
					if bar.Time.Unix() > lastMBar.Time.Unix() {
						lastMBar = bar
					}
				}

				if lastMBar.Time.Format("01-02-2006") != lastBar.Time.Format("01-02-2006") {
					return
				}

				// Get last bar from redis
				lastDBar, err := u.barRedisRepo.GetObj(ctx, u.genRedisLastBar(ticker.Symbol, resolutionD))
				if err != nil {
					errCh <- err
				}

				if lastMBar.Time.Format("01-02-2006") != lastDBar.Time.Format("01-02-2006") {
					mustSyncDTicker <- ticker
					return
				}

				newLastDBar := &models.Bar{
					Symbol:   lastDBar.Symbol,
					Exchange: lastDBar.Exchange,
					Time:     lastDBar.Time,
					Open:     lastDBar.Open,
					High:     lastDBar.High,
					Low:      lastDBar.Low,
					Close:    lastDBar.Close,
					Volume:   lastDBar.Volume,
				}
				hasNew := false

				for _, bar := range bars {
					// TODO: Convert to symbol timezone before compare
					if bar.Time.Format("01-02-2006") == newLastDBar.Time.Format("01-02-2006") {
						if bar.High > newLastDBar.High {
							newLastDBar.High = bar.High
							hasNew = true
						}

						if bar.Low < newLastDBar.Low {
							newLastDBar.Low = bar.Low
							hasNew = true
						}

						if lastBar.Time.Unix() < bar.Time.Unix() {
							newLastDBar.Volume += bar.Volume
							hasNew = true
						}
					}
				}

				if hasNew {
					newLastDBarCh <- newLastDBar
				}
			}(barsCh, newLastDBarCh, mustSyncDTicker, errSendCh, ticker)
		}
		sendWg.Wait()
		close(barsCh)
		close(newLastDBarCh)
		close(mustSyncDTicker)
	}(barsCh, newLastDBarCh, mustSyncDTicker, errSendCh)

	// Save
	doneCh := make(chan bool)
	errReceiveCh := make(chan error)

	go func() {
		var wg sync.WaitGroup

		wg.Add(1)
		go func(barsCh <-chan []*models.Bar, errCh chan<- error) {
			defer wg.Done()

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
		}(barsCh, errReceiveCh)

		wg.Add(1)
		go func(newLastDBarCh <-chan *models.Bar, errCh chan<- error) {
			defer wg.Done()

			var receiveWg sync.WaitGroup

			for newLastDBar := range newLastDBarCh {
				receiveWg.Add(1)
				go func(newLastDBar *models.Bar) {
					defer receiveWg.Done()

					err = u.Insert(ctx, resolutionD, newLastDBar, false)
					if err != nil {
						errCh <- err
					}
				}(newLastDBar)
			}

			receiveWg.Wait()
		}(newLastDBarCh, errReceiveCh)

		wg.Add(1)
		go func(mustSyncDTicker <-chan *models.Ticker, errCh chan<- error) {
			defer wg.Done()

			var tickers []*models.Ticker

			for ticker := range mustSyncDTicker {
				tickers = append(tickers, ticker)
			}

			err := u.syncDAllSymbol(ctx, tickerDownloadBatchSize, tickerInsertBatchSize, barInsertBatchSize, tickers)
			if err != nil {
				errCh <- err
			}
		}(mustSyncDTicker, errReceiveCh)

		wg.Wait()
		doneCh <- true
	}()

	select {
	case err := <-errSendCh:
		return err
	case err := <-errReceiveCh:
		return err
	case <-doneCh:
		return nil
	}
}

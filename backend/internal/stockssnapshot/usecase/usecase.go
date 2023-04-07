package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/stockssnapshot"
	"github.com/hiennguyen9874/stockk-go/internal/tickers"
	"github.com/hiennguyen9874/stockk-go/pkg/crawlers"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

type stockSnapshotUseCase struct {
	tickerUC               tickers.TickerUseCaseI
	stockSnapshotRedisRepo stockssnapshot.StockSnapshotRedisRepository
	crawler                crawlers.Crawler
	cfg                    *config.Config //nolint:unused
	logger                 logger.Logger  //nolint:unused
}

func CreateTickerUseCaseI(
	tickerUC tickers.TickerUseCaseI,
	stockSnapshotRedisRepo stockssnapshot.StockSnapshotRedisRepository,
	cfg *config.Config,
	logger logger.Logger,
) stockssnapshot.StockSnapshotUseCaseI {
	return &stockSnapshotUseCase{
		tickerUC:               tickerUC,
		stockSnapshotRedisRepo: stockSnapshotRedisRepo,
		crawler:                crawlers.NewCrawler(cfg, logger),
	}
}

func (u *stockSnapshotUseCase) GenerateRedisStockSnapshotKey(symbol string) string {
	return fmt.Sprintf("StockSnapshot:%v", symbol)
}

func (u *stockSnapshotUseCase) CrawlAllStocksSnapshot(ctx context.Context) error {
	tickers, err := u.tickerUC.GetAllActive(ctx, true)
	if err != nil {
		return err
	}

	var symbols []string
	for _, ticker := range tickers {
		symbols = append(symbols, ticker.Symbol)
	}

	stocksSnapshot, err := u.crawler.CrawlStockSnapshot(ctx, symbols)
	if err != nil {
		return err
	}

	doneCh := make(chan bool)
	errCh := make(chan error)

	go func() {
		var wg sync.WaitGroup

		for _, stockSnapshot := range stocksSnapshot {
			wg.Add(1)

			go func(stockSnapshot crawlers.StockSnapshot) {
				defer wg.Done()

				err := u.stockSnapshotRedisRepo.CreateObj(ctx, u.GenerateRedisStockSnapshotKey(stockSnapshot.Ticker), &models.StockSnapshot{
					Ticker:      stockSnapshot.Ticker,
					RefPrice:    stockSnapshot.RefPrice,
					CeilPrice:   stockSnapshot.CeilPrice,
					FloorPrice:  stockSnapshot.FloorPrice,
					TltVol:      stockSnapshot.TltVol,
					TltVal:      stockSnapshot.TltVal,
					PriceB3:     stockSnapshot.PriceB3,
					PriceB2:     stockSnapshot.PriceB2,
					PriceB1:     stockSnapshot.PriceB1,
					VolB3:       stockSnapshot.VolB3,
					VolB2:       stockSnapshot.VolB2,
					VolB1:       stockSnapshot.VolB1,
					Price:       stockSnapshot.Price,
					Vol:         stockSnapshot.Vol,
					PriceS3:     stockSnapshot.PriceS3,
					PriceS2:     stockSnapshot.PriceS2,
					PriceS1:     stockSnapshot.PriceS1,
					VolS3:       stockSnapshot.VolS3,
					VolS2:       stockSnapshot.VolS2,
					VolS1:       stockSnapshot.VolS1,
					High:        stockSnapshot.High,
					Low:         stockSnapshot.Low,
					BuyForeign:  stockSnapshot.BuyForeign,
					SellForeign: stockSnapshot.SellForeign,
				}, -1)

				if err != nil {
					errCh <- err
				}
			}(stockSnapshot)
		}
		wg.Wait()

		doneCh <- true
	}()

	select {
	case err := <-errCh:
		return err
	case <-doneCh:
		return nil
	}
}

func (u *stockSnapshotUseCase) GetStockSnapshotBySymbol(ctx context.Context, symbol string) (*models.StockSnapshot, error) {
	return u.stockSnapshotRedisRepo.GetObj(ctx, u.GenerateRedisStockSnapshotKey(symbol))
}

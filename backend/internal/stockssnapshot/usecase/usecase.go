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
	crawler                crawlers.RestCrawler
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

				err := u.stockSnapshotRedisRepo.CreateObj(ctx, u.GenerateRedisStockSnapshotKey(stockSnapshot.Code), &models.StockSnapshot{
					Ticker:          stockSnapshot.Code,
					BasicPrice:      stockSnapshot.BasicPrice,
					CeilingPrice:    stockSnapshot.CeilingPrice,
					FloorPrice:      stockSnapshot.FloorPrice,
					AccumulatedVol:  stockSnapshot.AccumulatedVol,
					AccumulatedVal:  stockSnapshot.AccumulatedVal,
					MatchPrice:      stockSnapshot.MatchPrice,
					MatchQtty:       stockSnapshot.MatchQtty,
					HighestPrice:    stockSnapshot.HighestPrice,
					LowestPrice:     stockSnapshot.LowestPrice,
					BuyForeignQtty:  stockSnapshot.BuyForeignQtty,
					SellForeignQtty: stockSnapshot.SellForeignQtty,
					ProjectOpen:     stockSnapshot.ProjectOpen,
					CurrentRoom:     stockSnapshot.CurrentRoom,
					FloorCode:       stockSnapshot.FloorCode,
					TotalRoom:       stockSnapshot.TotalRoom,
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

func (u *stockSnapshotUseCase) UpdateStockSnapshotBySymbol(ctx context.Context, symbol string, values map[string]interface{}) error {
	stockSnapshot, err := u.stockSnapshotRedisRepo.GetObj(ctx, u.GenerateRedisStockSnapshotKey(symbol))
	if err != nil {
		return err
	}

	if stockSnapshot == nil {
		stockSnapshot = &models.StockSnapshot{
			Ticker: symbol,
		}
	}

	if value, ok := values["BasicPrice"]; ok {
		if valueFloat, ok := value.(float32); ok {
			stockSnapshot.BasicPrice = valueFloat
		}
	}
	if value, ok := values["CeilingPrice"]; ok {
		if valueFloat, ok := value.(float32); ok {
			stockSnapshot.CeilingPrice = valueFloat
		}
	}
	if value, ok := values["FloorPrice"]; ok {
		if valueFloat, ok := value.(float32); ok {
			stockSnapshot.FloorPrice = valueFloat
		}
	}
	if value, ok := values["AccumulatedVol"]; ok {
		if valueFloat, ok := value.(float32); ok {
			stockSnapshot.AccumulatedVol = valueFloat
		}
	}
	if value, ok := values["AccumulatedVal"]; ok {
		if valueFloat, ok := value.(float32); ok {
			stockSnapshot.AccumulatedVal = valueFloat
		}
	}
	if value, ok := values["MatchPrice"]; ok {
		if valueFloat, ok := value.(float32); ok {
			stockSnapshot.MatchPrice = valueFloat
		}
	}
	if value, ok := values["MatchQtty"]; ok {
		if valueFloat, ok := value.(float32); ok {
			stockSnapshot.MatchQtty = valueFloat
		}
	}
	if value, ok := values["HighestPrice"]; ok {
		if valueFloat, ok := value.(float32); ok {
			stockSnapshot.HighestPrice = valueFloat
		}
	}
	if value, ok := values["LowestPrice"]; ok {
		if valueFloat, ok := value.(float32); ok {
			stockSnapshot.LowestPrice = valueFloat
		}
	}
	if value, ok := values["BuyForeignQtty"]; ok {
		if valueFloat, ok := value.(float32); ok {
			stockSnapshot.BuyForeignQtty = valueFloat
		}
	}
	if value, ok := values["SellForeignQtty"]; ok {
		if valueFloat, ok := value.(float32); ok {
			stockSnapshot.SellForeignQtty = valueFloat
		}
	}
	if value, ok := values["ProjectOpen"]; ok {
		if valueFloat, ok := value.(float32); ok {
			stockSnapshot.ProjectOpen = valueFloat
		}
	}
	if value, ok := values["CurrentRoom"]; ok {
		if valueFloat, ok := value.(float32); ok {
			stockSnapshot.CurrentRoom = valueFloat
		}
	}
	if value, ok := values["FloorCode"]; ok {
		if valueStr, ok := value.(string); ok {
			stockSnapshot.FloorCode = valueStr
		}
	}
	if value, ok := values["TotalRoom"]; ok {
		if valueFloat, ok := value.(float32); ok {
			stockSnapshot.TotalRoom = valueFloat
		}
	}

	return u.stockSnapshotRedisRepo.CreateObj(ctx, u.GenerateRedisStockSnapshotKey(stockSnapshot.Ticker), stockSnapshot, -1)
}

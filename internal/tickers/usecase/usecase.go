package usecase

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/tickers"
	"github.com/hiennguyen9874/stockk-go/internal/usecase"
	"github.com/hiennguyen9874/stockk-go/pkg/crawlers"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

const (
	batchSizeSaveTickers = 500
)

type tickerUseCase struct {
	usecase.UseCase[models.Ticker]
	pgRepo  tickers.TickerPgRepository
	crawler crawlers.Crawler
}

func CreateTickerUseCaseI(
	pgRepo tickers.TickerPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) tickers.TickerUseCaseI {
	return &tickerUseCase{
		UseCase: usecase.CreateUseCase[models.Ticker](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
		crawler: crawlers.NewCrawler(cfg),
	}
}

func (u *tickerUseCase) GetBySymbol(ctx context.Context, symbol string) (*models.Ticker, error) {
	return u.pgRepo.GetBySymbol(ctx, symbol)
}

func (u *tickerUseCase) CrawlAllStockTicker(ctx context.Context) ([]*models.Ticker, error) {
	tickers, err := u.crawler.VNDCrawlStockSymbols()
	if err != nil {
		return nil, err
	}

	inDBTickers, err := u.pgRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	keyTickers := make(map[string]bool)
	for _, ticker := range inDBTickers {
		keyTickers[ticker.Symbol] = true
	}

	var mustCreateTickers []*models.Ticker
	for _, ticker := range tickers {
		if _, ok := keyTickers[ticker.Symbol]; !ok {
			mustCreateTickers = append(mustCreateTickers, &models.Ticker{
				Symbol:    ticker.Symbol,
				Exchange:  ticker.Exchange,
				FullName:  ticker.FullName,
				ShortName: ticker.ShortName,
				Type:      ticker.Type,
				IsActive:  false,
			})
		}
	}

	savedTickers, err := u.pgRepo.CreateMulti(ctx, mustCreateTickers, batchSizeSaveTickers)
	if err != nil {
		return nil, err
	}

	return savedTickers, nil
}

func (u *tickerUseCase) UpdateIsActiveBySymbol(ctx context.Context, symbol string, isActive bool) (*models.Ticker, error) {
	ticker, err := u.pgRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	u.Logger.Info(isActive)
	updatedTicker, err := u.pgRepo.UpdateIsActive(ctx, ticker, isActive)
	if err != nil {
		return nil, err
	}
	return updatedTicker, nil
}

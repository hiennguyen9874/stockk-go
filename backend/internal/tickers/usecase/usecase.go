package usecase

import (
	"context"
	"fmt"

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
	tickerPgRepo    tickers.TickerPgRepository
	tickerRedisRepo tickers.TickerRedisRepository
	crawler         crawlers.Crawler
}

func CreateTickerUseCaseI(
	tickerPgRepo tickers.TickerPgRepository,
	tickerRedisRepo tickers.TickerRedisRepository,
	cfg *config.Config,
	logger logger.Logger,
) tickers.TickerUseCaseI {
	return &tickerUseCase{
		UseCase:         usecase.CreateUseCase[models.Ticker](tickerPgRepo, cfg, logger),
		tickerPgRepo:    tickerPgRepo,
		tickerRedisRepo: tickerRedisRepo,
		crawler:         crawlers.NewCrawler(cfg, logger),
	}
}

func (u *tickerUseCase) GenerateRedisTickerKey(symbol string) string {
	return fmt.Sprintf("%v:%v", models.Ticker{}.TableName(), symbol)
}

func (u *tickerUseCase) GenerateRedisAllTickerActive(isActive bool) string {
	return fmt.Sprintf("Ticker:AllTicker:%v", isActive)
}

func (u *tickerUseCase) GetBySymbol(ctx context.Context, symbol string) (*models.Ticker, error) {
	cachedTicker, err := u.tickerRedisRepo.GetObj(ctx, u.GenerateRedisTickerKey(symbol))
	if err != nil {
		return nil, err
	}

	if cachedTicker != nil {
		return cachedTicker, nil
	}

	ticker, err := u.tickerPgRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	if err = u.tickerRedisRepo.CreateObj(ctx, u.GenerateRedisTickerKey(symbol), ticker, 3600); err != nil {
		return nil, err
	}

	return ticker, nil
}

func (u *tickerUseCase) CrawlAllStockTicker(ctx context.Context) ([]*models.Ticker, error) {
	tickers, err := u.crawler.CrawlStockSymbols(ctx)
	if err != nil {
		return nil, err
	}

	inDBTickers, err := u.tickerPgRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	keyTickers := make(map[string]bool)
	for _, ticker := range inDBTickers {
		keyTickers[ticker.Symbol] = true
	}

	defaultActiveTickers := make(map[string]bool, len(u.Cfg.Crawler.DefaultActive))
	for _, ticker := range u.Cfg.Crawler.DefaultActive {
		defaultActiveTickers[ticker] = true
	}

	var mustCreateTickers []*models.Ticker
	for _, ticker := range tickers {
		if _, ok := keyTickers[ticker.Symbol]; !ok {
			_, active := defaultActiveTickers[ticker.Symbol]

			mustCreateTickers = append(mustCreateTickers, &models.Ticker{
				Symbol:    ticker.Symbol,
				Exchange:  ticker.Exchange,
				FullName:  ticker.FullName,
				ShortName: ticker.ShortName,
				Type:      ticker.Type,
				IsActive:  active,
			})
		}
	}

	savedTickers, err := u.tickerPgRepo.CreateMulti(ctx, mustCreateTickers, batchSizeSaveTickers)
	if err != nil {
		return nil, err
	}

	if len(savedTickers) > 0 {
		if err = u.tickerRedisRepo.Delete(ctx, u.GenerateRedisAllTickerActive(true)); err != nil {
			return nil, err
		}
		if err = u.tickerRedisRepo.Delete(ctx, u.GenerateRedisAllTickerActive(false)); err != nil {
			return nil, err
		}
	}

	return savedTickers, nil
}

func (u *tickerUseCase) UpdateIsActiveBySymbol(ctx context.Context, symbol string, isActive bool) (*models.Ticker, error) {
	ticker, err := u.tickerPgRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	u.Logger.Info(isActive)
	updatedTicker, err := u.tickerPgRepo.UpdateIsActive(ctx, ticker, isActive)
	if err != nil {
		return nil, err
	}

	if err = u.tickerRedisRepo.Delete(ctx, u.GenerateRedisTickerKey(symbol)); err != nil {
		return nil, err
	}

	if err = u.tickerRedisRepo.Delete(ctx, u.GenerateRedisAllTickerActive(true)); err != nil {
		return nil, err
	}

	if err = u.tickerRedisRepo.Delete(ctx, u.GenerateRedisAllTickerActive(false)); err != nil {
		return nil, err
	}

	return updatedTicker, nil
}

func (u *tickerUseCase) GetAllActive(ctx context.Context, isActive bool) ([]*models.Ticker, error) {
	cachedTickers, err := u.tickerRedisRepo.GetObjs(ctx, u.GenerateRedisAllTickerActive(true))
	if err != nil {
		return nil, err
	}

	if cachedTickers != nil {
		return cachedTickers, nil
	}

	tickers, err := u.tickerPgRepo.GetAllActive(ctx, isActive)
	if err != nil {
		return nil, err
	}

	if err = u.tickerRedisRepo.CreateObjs(ctx, u.GenerateRedisAllTickerActive(true), tickers, 3600); err != nil {
		return nil, err
	}

	return tickers, nil
}

func (u *tickerUseCase) SearchBySymbol(ctx context.Context, symbol string, limit int, exchange string, isActive bool) ([]*models.Ticker, error) {
	return u.tickerPgRepo.SearchBySymbol(ctx, symbol, limit, exchange, isActive)
}

package usecase

import (
	"github.com/hiennguyen9874/go-boilerplate/config"
	"github.com/hiennguyen9874/go-boilerplate/internal/models"
	"github.com/hiennguyen9874/go-boilerplate/internal/tickers"
	"github.com/hiennguyen9874/go-boilerplate/internal/usecase"
	"github.com/hiennguyen9874/go-boilerplate/pkg/logger"
)

type tickerUseCase struct {
	usecase.UseCase[models.Ticker]
	pgRepo tickers.TickerPgRepository
}

func CreateTickerUseCaseI(
	pgRepo tickers.TickerPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) tickers.TickerUseCaseI {
	return &tickerUseCase{
		UseCase: usecase.CreateUseCase[models.Ticker](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

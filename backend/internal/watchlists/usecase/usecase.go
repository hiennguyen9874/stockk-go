package usecase

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/usecase"
	"github.com/hiennguyen9874/stockk-go/internal/watchlists"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

type watchListUseCase struct {
	usecase.UseCase[models.WatchList]
	pgRepo watchlists.WatchListPgRepository
}

func CreateWatchListUseCaseI(
	pgRepo watchlists.WatchListPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) watchlists.WatchListUseCaseI {
	return &watchListUseCase{
		UseCase: usecase.CreateUseCase[models.WatchList](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

func (u *watchListUseCase) GetMultiByOwnerId(ctx context.Context, ownerId uint, limit, offset int) ([]*models.WatchList, error) {
	return u.pgRepo.GetMultiByOwnerId(ctx, ownerId, limit, offset)
}

func (u *watchListUseCase) CreateWithOwner(ctx context.Context, ownerId uint, exp *models.WatchList) (*models.WatchList, error) {
	return u.pgRepo.CreateWithOwner(ctx, ownerId, exp)
}

func (u *watchListUseCase) DeleteWithoutGet(ctx context.Context, id uint) error {
	return u.pgRepo.DeleteWithoutGet(ctx, id)
}

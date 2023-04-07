package usecase

import (
	"context"
	"errors"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/clients"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/usecase"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

type clientUseCase struct {
	usecase.UseCase[models.Client]
	pgRepo clients.ClientPgRepository
}

func CreateClientUseCaseI(
	pgRepo clients.ClientPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) clients.ClientUseCaseI {
	return &clientUseCase{
		UseCase: usecase.CreateUseCase[models.Client](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

func (u *clientUseCase) GetMultiByOwnerId(ctx context.Context, ownerId uint, limit, offset int) ([]*models.Client, error) {
	return u.pgRepo.GetMultiByOwnerId(ctx, ownerId, limit, offset)
}

func (u *clientUseCase) CreateWithOwner(ctx context.Context, ownerId uint, exp *models.Client) (*models.Client, error) {
	count, err := u.pgRepo.CountWithOwner(ctx, ownerId)
	if err != nil {
		return nil, err
	}

	if count >= 8 {
		return nil, errors.New("can not create more than 8 chart")
	}

	return u.pgRepo.CreateWithOwner(ctx, ownerId, exp)
}

func (u *clientUseCase) DeleteWithoutGet(ctx context.Context, id uint) error {
	ticker, err := u.pgRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	count, err := u.pgRepo.CountWithOwner(ctx, ticker.OwnerId)
	if err != nil {
		return err
	}

	if count <= 1 {
		return errors.New("can not remove when remain only one chart")
	}

	return u.pgRepo.DeleteWithoutGet(ctx, id)
}

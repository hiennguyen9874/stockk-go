package usecase

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/charts"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/usecase"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

type chartUseCase struct {
	usecase.UseCase[models.Chart]
	chartPgRepo charts.ChartPgRepository
}

func CreateChartUseCaseI(
	chartPgRepo charts.ChartPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) charts.ChartUseCaseI {
	return &chartUseCase{
		UseCase:     usecase.CreateUseCase[models.Chart](chartPgRepo, cfg, logger),
		chartPgRepo: chartPgRepo,
	}
}

func (u *chartUseCase) GetAllByOwner(ctx context.Context, ownerSource string, ownerId string) ([]*models.Chart, error) {
	return u.chartPgRepo.GetAllByOwner(ctx, ownerSource, ownerId)
}

func (u *chartUseCase) GetByIdOwner(ctx context.Context, id uint, ownerSource string, ownerId string) (*models.Chart, error) {
	return u.chartPgRepo.GetByIdOwner(ctx, id, ownerSource, ownerId)
}

func (u *chartUseCase) DeleteByIdOwner(ctx context.Context, id uint, ownerSource string, ownerId string) (*models.Chart, error) {
	return u.chartPgRepo.DeleteByIdOwner(ctx, id, ownerSource, ownerId)
}

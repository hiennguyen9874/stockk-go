package usecase

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

type UseCase[M any] struct {
	Cfg    *config.Config
	pgRepo internal.PgRepository[M]
	Logger logger.Logger
}

func CreateUseCase[M any](repo internal.PgRepository[M], cfg *config.Config, logger logger.Logger) UseCase[M] {
	return UseCase[M]{
		pgRepo: repo,
		Cfg:    cfg,
		Logger: logger,
	}
}

func CreateUseCaseI[M any](repo internal.PgRepository[M], cfg *config.Config, logger logger.Logger) internal.UseCaseI[M] {
	return &UseCase[M]{
		pgRepo: repo,
		Cfg:    cfg,
		Logger: logger,
	}
}

func (u *UseCase[M]) Create(ctx context.Context, exp *M) (*M, error) {
	return u.pgRepo.Create(ctx, exp)
}

func (u *UseCase[M]) Get(ctx context.Context, id uint) (*M, error) {
	return u.pgRepo.Get(ctx, id)
}

func (u *UseCase[M]) GetMulti(ctx context.Context, limit int, offset int) ([]*M, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return u.pgRepo.GetMulti(ctx, limit, offset)
}

func (u *UseCase[M]) Delete(ctx context.Context, id uint) (*M, error) {
	return u.pgRepo.Delete(ctx, id)
}

func (u *UseCase[M]) Update(ctx context.Context, id uint, values map[string]interface{}) (*M, error) {
	obj, err := u.Get(ctx, id)
	if err != nil || obj == nil {
		return nil, err
	}

	return u.pgRepo.Update(ctx, obj, values)
}

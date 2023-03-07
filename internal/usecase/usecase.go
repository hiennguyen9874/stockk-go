package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/hiennguyen9874/stockk-go/internal"
)

type UseCase[M any] struct {
	// TODO: Config
	// TODO: Logger
	pgRepo internal.PgRepository[M]
}

func CreateUseCase[M any](repo internal.PgRepository[M]) UseCase[M] {
	return UseCase[M]{
		pgRepo: repo,
	}
}

func (u *UseCase[M]) Create(ctx context.Context, exp *M) (*M, error) {
	return u.pgRepo.Create(ctx, exp)
}

func (u *UseCase[M]) Get(ctx context.Context, id uuid.UUID) (*M, error) {
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

func (u *UseCase[M]) Delete(ctx context.Context, id uuid.UUID) (*M, error) {
	return u.pgRepo.Delete(ctx, id)
}

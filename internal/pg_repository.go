package internal

import (
	"context"

	"github.com/google/uuid"
)

type PgRepository[M any] interface {
	Create(ctx context.Context, exp *M) (res *M, err error)
	Get(ctx context.Context, id uuid.UUID) (res *M, err error)
	GetMulti(ctx context.Context, limit, offset int) (res []*M, err error)
	Delete(ctx context.Context, id uuid.UUID) (res *M, err error)
	Update(ctx context.Context, exp *M, values map[string]interface{}) (res *M, err error)
}

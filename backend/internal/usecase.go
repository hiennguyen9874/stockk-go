package internal

import (
	"context"
)

type UseCaseI[M any] interface {
	Create(ctx context.Context, exp *M) (*M, error)
	Get(ctx context.Context, id uint) (*M, error)
	GetMulti(ctx context.Context, limit int, offset int) ([]*M, error)
	Delete(ctx context.Context, id uint) (*M, error)
	Update(ctx context.Context, id uint, values map[string]interface{}) (*M, error)
}

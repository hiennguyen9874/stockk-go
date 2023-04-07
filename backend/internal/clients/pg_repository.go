package clients

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type ClientPgRepository interface {
	internal.PgRepository[models.Client]
	GetMultiByOwnerId(ctx context.Context, ownerId uint, limit, offset int) ([]*models.Client, error)
	CreateWithOwner(ctx context.Context, ownerId uint, exp *models.Client) (*models.Client, error)
	DeleteWithoutGet(ctx context.Context, id uint) error
	CountWithOwner(ctx context.Context, ownerId uint) (int64, error)
}

package watchlists

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type WatchListUseCaseI interface {
	internal.UseCaseI[models.WatchList]
	GetMultiByOwnerId(ctx context.Context, ownerId uint, limit, offset int) ([]*models.WatchList, error)
	CreateWithOwner(ctx context.Context, ownerId uint, exp *models.WatchList) (*models.WatchList, error)
	DeleteWithoutGet(ctx context.Context, id uint) error
}

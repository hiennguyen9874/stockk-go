package charts

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type ChartPgRepository interface {
	internal.PgRepository[models.Chart]
	GetAllByOwner(ctx context.Context, ownerSource string, ownerId string) ([]*models.Chart, error)
	GetByIdOwner(ctx context.Context, id uint, ownerSource string, ownerId string) (*models.Chart, error)
	DeleteByIdOwner(ctx context.Context, id uint, ownerSource string, ownerId string) (*models.Chart, error)
}

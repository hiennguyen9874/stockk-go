package studytemplates

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type StudyTemplatePgRepository interface {
	internal.PgRepository[models.StudyTemplate]
	GetByOwnerName(ctx context.Context, ownerSource string, ownerId string, name string) (*models.StudyTemplate, error)
	GetOrCreateWithOwnerName(ctx context.Context, ownerSource string, ownerId string, name string, content string) (*models.StudyTemplate, bool, error)
	GetAllByOwner(ctx context.Context, ownerSource string, ownerId string) ([]*models.StudyTemplate, error)
	DeleteByOwnerName(ctx context.Context, ownerSource string, ownerId string, name string) (*models.StudyTemplate, error)
}

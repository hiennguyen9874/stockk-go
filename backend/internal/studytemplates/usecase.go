package studytemplates

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type StudyTemplateUseCaseI interface {
	internal.UseCaseI[models.StudyTemplate]
	GetByOwnerName(ctx context.Context, ownerSource string, ownerId string, name string) (*models.StudyTemplate, error)
	GetOrCreateWithOwnerName(ctx context.Context, ownerSource string, ownerId string, name string, content string) (*models.StudyTemplate, bool, error)
	GetAllByOwner(ctx context.Context, ownerSource string, ownerId string) ([]*models.StudyTemplate, error)
	CreateOrUpdateWithOwnerName(ctx context.Context, ownerSource string, ownerId string, name string, content string) (*models.StudyTemplate, bool, bool, error)
	DeleteByOwnerName(ctx context.Context, ownerSource string, ownerId string, name string) (*models.StudyTemplate, error)
}

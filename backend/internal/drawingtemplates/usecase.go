package drawingtemplates

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type DrawingTemplateUseCaseI interface {
	internal.UseCaseI[models.DrawingTemplate]
	GetByOwnerNameTool(ctx context.Context, ownerSource string, ownerId string, name string, tool string) (*models.DrawingTemplate, error)
	GetOrCreateWithOwnerNameTool(ctx context.Context, ownerSource string, ownerId string, name string, tool string, content string) (*models.DrawingTemplate, bool, error)
	GetAllByOwnerTool(ctx context.Context, ownerSource string, ownerId string, tool string) ([]*models.DrawingTemplate, error)
	DeleteByOwnerNameTool(ctx context.Context, ownerSource string, ownerId string, name string, tool string) (*models.DrawingTemplate, error)
	CreateOrUpdateWithOwnerNameTool(ctx context.Context, ownerSource string, ownerId string, name string, tool string, content string) (*models.DrawingTemplate, bool, bool, error)
}

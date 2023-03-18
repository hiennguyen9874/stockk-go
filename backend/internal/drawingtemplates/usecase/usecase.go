package usecase

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/drawingtemplates"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/usecase"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

type drawingTemplateUseCase struct {
	usecase.UseCase[models.DrawingTemplate]
	drawingTemplatePgRepo drawingtemplates.DrawingTemplatePgRepository
}

func DrawingTemplateUseCaseI(
	drawingTemplatePgRepo drawingtemplates.DrawingTemplatePgRepository,
	cfg *config.Config,
	logger logger.Logger,
) drawingtemplates.DrawingTemplateUseCaseI {
	return &drawingTemplateUseCase{
		UseCase:               usecase.CreateUseCase[models.DrawingTemplate](drawingTemplatePgRepo, cfg, logger),
		drawingTemplatePgRepo: drawingTemplatePgRepo,
	}
}

func (u *drawingTemplateUseCase) GetByOwnerNameTool(ctx context.Context, ownerSource string, ownerId string, name string, tool string) (*models.DrawingTemplate, error) {
	return u.drawingTemplatePgRepo.GetByOwnerNameTool(ctx, ownerSource, ownerId, name, tool)
}

func (u *drawingTemplateUseCase) GetOrCreateWithOwnerNameTool(ctx context.Context, ownerSource string, ownerId string, name string, tool string, content string) (*models.DrawingTemplate, bool, error) {
	return u.drawingTemplatePgRepo.GetOrCreateWithOwnerNameTool(ctx, ownerSource, ownerId, name, tool, content)
}

func (u *drawingTemplateUseCase) GetAllByOwnerTool(ctx context.Context, ownerSource string, ownerId string, tool string) ([]*models.DrawingTemplate, error) {
	return u.drawingTemplatePgRepo.GetAllByOwnerTool(ctx, ownerSource, ownerId, tool)
}

func (u *drawingTemplateUseCase) DeleteByOwnerNameTool(ctx context.Context, ownerSource string, ownerId string, name string, tool string) (*models.DrawingTemplate, error) {
	return u.drawingTemplatePgRepo.DeleteByOwnerNameTool(ctx, ownerSource, ownerId, name, tool)
}

func (u *drawingTemplateUseCase) CreateOrUpdateWithOwnerNameTool(ctx context.Context, ownerSource string, ownerId string, name string, tool string, content string) (*models.DrawingTemplate, bool, bool, error) {
	drawingTemplate, created, err := u.drawingTemplatePgRepo.GetOrCreateWithOwnerNameTool(ctx, ownerSource, ownerId, name, tool, content)

	if err != nil {
		return nil, false, false, err
	}

	if created {
		return drawingTemplate, true, false, nil
	}

	values := make(map[string]interface{})
	values["content"] = content

	updatedDrawingTemplate, err := u.drawingTemplatePgRepo.Update(ctx, drawingTemplate, values)
	if err != nil {
		return nil, false, false, err
	}

	return updatedDrawingTemplate, false, true, nil
}

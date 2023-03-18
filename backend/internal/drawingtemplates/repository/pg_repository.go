package repository

import (
	"context"
	"errors"

	"github.com/hiennguyen9874/stockk-go/internal/drawingtemplates"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/repository"
	"gorm.io/gorm"
)

type DrawingTemplatePgRepo struct {
	repository.PgRepo[models.DrawingTemplate]
}

func CreateDrawingTemplatePgRepository(db *gorm.DB) drawingtemplates.DrawingTemplatePgRepository {
	return &DrawingTemplatePgRepo{
		PgRepo: repository.CreatePgRepo[models.DrawingTemplate](db),
	}
}

func (r *DrawingTemplatePgRepo) GetByOwnerNameTool(ctx context.Context, ownerSource string, ownerId string, name string, tool string) (*models.DrawingTemplate, error) {
	var obj *models.DrawingTemplate
	if result := r.DB.WithContext(ctx).First(&obj, "name = ? AND tool = ? AND owner_source = ? AND owner_id = ?", name, tool, ownerSource, ownerId); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *DrawingTemplatePgRepo) GetOrCreateWithOwnerNameTool(ctx context.Context, ownerSource string, ownerId string, name string, tool string, content string) (*models.DrawingTemplate, bool, error) {
	var obj *models.DrawingTemplate
	if result := r.DB.WithContext(ctx).First(&obj, "name = ? AND tool = ? AND owner_source = ? AND owner_id = ?", name, tool, ownerSource, ownerId); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			newDrawingTemplate := &models.DrawingTemplate{
				OwnerSource: ownerSource,
				OwnerId:     ownerId,
				Name:        name,
				Content:     content,
			}

			if result := r.DB.WithContext(ctx).Create(newDrawingTemplate); result.Error != nil {
				return nil, false, result.Error
			}
			return newDrawingTemplate, true, nil
		}
		return nil, false, result.Error
	}
	return obj, false, nil
}

func (r *DrawingTemplatePgRepo) GetAllByOwnerTool(ctx context.Context, ownerSource string, ownerId string, tool string) ([]*models.DrawingTemplate, error) {
	var objs []*models.DrawingTemplate
	r.DB.WithContext(ctx).Where("tool = ? AND owner_source = ? AND owner_id = ?", tool, ownerSource, ownerId).Find(&objs)
	return objs, nil
}

func (r *DrawingTemplatePgRepo) DeleteByOwnerNameTool(ctx context.Context, ownerSource string, ownerId string, name string, tool string) (*models.DrawingTemplate, error) {
	obj, err := r.GetByOwnerNameTool(ctx, ownerSource, ownerId, name, tool)

	if err != nil {
		return nil, err
	}

	if result := r.DB.WithContext(ctx).Delete(&obj, "name = ? AND tool = ? AND owner_source = ? AND owner_id = ?", name, tool, ownerSource, ownerId); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

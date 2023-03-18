package repository

import (
	"context"
	"errors"

	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/repository"
	"github.com/hiennguyen9874/stockk-go/internal/studytemplates"
	"gorm.io/gorm"
)

type StudyTemplatePgRepo struct {
	repository.PgRepo[models.StudyTemplate]
}

func CreateStudyTemplatePgRepository(db *gorm.DB) studytemplates.StudyTemplatePgRepository {
	return &StudyTemplatePgRepo{
		PgRepo: repository.CreatePgRepo[models.StudyTemplate](db),
	}
}

func (r *StudyTemplatePgRepo) GetByOwnerName(ctx context.Context, ownerSource string, ownerId string, name string) (*models.StudyTemplate, error) {
	var obj *models.StudyTemplate
	if result := r.DB.WithContext(ctx).First(&obj, "name = ? AND owner_source = ? AND owner_id = ?", name, ownerSource, ownerId); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *StudyTemplatePgRepo) GetOrCreateWithOwnerName(ctx context.Context, ownerSource string, ownerId string, name string, content string) (*models.StudyTemplate, bool, error) {
	var obj *models.StudyTemplate
	if result := r.DB.WithContext(ctx).First(&obj, "name = ? AND owner_source = ? AND owner_id = ?", name, ownerSource, ownerId); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			newStudyTemplate := &models.StudyTemplate{
				OwnerSource: ownerSource,
				OwnerId:     ownerId,
				Name:        name,
				Content:     content,
			}

			if result := r.DB.WithContext(ctx).Create(newStudyTemplate); result.Error != nil {
				return nil, false, result.Error
			}
			return newStudyTemplate, true, nil
		}
		return nil, false, result.Error
	}
	return obj, false, nil
}

func (r *StudyTemplatePgRepo) GetAllByOwner(ctx context.Context, ownerSource string, ownerId string) ([]*models.StudyTemplate, error) {
	var objs []*models.StudyTemplate
	r.DB.WithContext(ctx).Where("owner_source = ? AND owner_id = ?", ownerSource, ownerId).Find(&objs)
	return objs, nil
}

func (r *StudyTemplatePgRepo) DeleteByOwnerName(ctx context.Context, ownerSource string, ownerId string, name string) (*models.StudyTemplate, error) {
	obj, err := r.GetByOwnerName(ctx, ownerSource, ownerId, name)

	if err != nil {
		return nil, err
	}

	if result := r.DB.WithContext(ctx).Delete(&obj, "name = ? AND owner_source = ? AND owner_id = ?", name, ownerSource, ownerId); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

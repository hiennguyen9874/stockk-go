package repository

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal/charts"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/repository"
	"gorm.io/gorm"
)

type ChartPgRepo struct {
	repository.PgRepo[models.Chart]
}

func CreateChartPgRepository(db *gorm.DB) charts.ChartPgRepository {
	return &ChartPgRepo{
		PgRepo: repository.CreatePgRepo[models.Chart](db),
	}
}

func (r *ChartPgRepo) GetAllByOwner(ctx context.Context, ownerSource string, ownerId string) ([]*models.Chart, error) {
	var objs []*models.Chart
	r.DB.WithContext(ctx).Where("owner_source = ? AND owner_id = ?", ownerSource, ownerId).Find(&objs)
	return objs, nil
}

func (r *ChartPgRepo) GetByIdOwner(ctx context.Context, id uint, ownerSource string, ownerId string) (*models.Chart, error) {
	var obj *models.Chart
	if result := r.DB.WithContext(ctx).First(&obj, "id = ? AND owner_source = ? AND owner_id = ?", id, ownerSource, ownerId); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *ChartPgRepo) DeleteByIdOwner(ctx context.Context, id uint, ownerSource string, ownerId string) (*models.Chart, error) {
	obj, err := r.GetByIdOwner(ctx, id, ownerSource, ownerId)

	if err != nil {
		return nil, err
	}

	if result := r.DB.WithContext(ctx).Delete(&obj, "id = ? AND owner_source = ? AND owner_id = ?", id, ownerSource, ownerId); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

package repository

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal/clients"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/repository"
	"gorm.io/gorm"
)

type ClientPgRepo struct {
	repository.PgRepo[models.Client]
}

func CreateClientPgRepository(db *gorm.DB) clients.ClientPgRepository {
	return &ClientPgRepo{
		PgRepo: repository.CreatePgRepo[models.Client](db),
	}
}

func (r *ClientPgRepo) GetMultiByOwnerId(ctx context.Context, ownerId uint, limit, offset int) ([]*models.Client, error) {
	var objs []*models.Client
	r.DB.WithContext(ctx).Where("owner_id = ?", ownerId).Order("created_at").Limit(limit).Offset(offset).Find(&objs)
	return objs, nil
}

func (r *ClientPgRepo) CreateWithOwner(ctx context.Context, ownerId uint, exp *models.Client) (*models.Client, error) {
	exp.OwnerId = ownerId
	if result := r.DB.WithContext(ctx).Create(exp); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *ClientPgRepo) DeleteWithoutGet(ctx context.Context, id uint) error {
	if result := r.DB.WithContext(ctx).Delete(&models.Client{}, "id = ?", id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ClientPgRepo) CountWithOwner(ctx context.Context, ownerId uint) (int64, error) {
	var count int64

	if result := r.DB.WithContext(ctx).Model(&models.Client{}).Where("owner_id = ?", ownerId).Count(&count); result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

package repository

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/repository"
	"github.com/hiennguyen9874/stockk-go/internal/watchlists"
	"gorm.io/gorm"
)

type WatchListPgRepo struct {
	repository.PgRepo[models.WatchList]
}

func CreateWatchListPgRepository(db *gorm.DB) watchlists.WatchListPgRepository {
	return &WatchListPgRepo{
		PgRepo: repository.CreatePgRepo[models.WatchList](db),
	}
}

func (r *WatchListPgRepo) GetMultiByOwnerId(ctx context.Context, ownerId uint, limit, offset int) ([]*models.WatchList, error) {
	var objs []*models.WatchList
	r.DB.WithContext(ctx).Where("owner_id = ?", ownerId).Order("created_at").Limit(limit).Offset(offset).Find(&objs)
	return objs, nil
}

func (r *WatchListPgRepo) CreateWithOwner(ctx context.Context, ownerId uint, exp *models.WatchList) (*models.WatchList, error) {
	exp.OwnerId = ownerId
	if result := r.DB.WithContext(ctx).Create(exp); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *WatchListPgRepo) DeleteWithoutGet(ctx context.Context, id uint) error {
	if result := r.DB.WithContext(ctx).Delete(&models.WatchList{}, "id = ?", id); result.Error != nil {
		return result.Error
	}
	return nil
}

package repository

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal"
	"gorm.io/gorm"
)

type PgRepo[M any] struct {
	DB *gorm.DB
}

func CreatePgRepo[M any](db *gorm.DB) PgRepo[M] {
	return PgRepo[M]{DB: db}
}

func CreatePgRepository[M any](db *gorm.DB) internal.PgRepository[M] {
	return &PgRepo[M]{DB: db}
}

func (r *PgRepo[M]) Get(ctx context.Context, id uint) (*M, error) {
	var obj *M
	if result := r.DB.WithContext(ctx).First(&obj, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *PgRepo[M]) GetMulti(ctx context.Context, limit, offset int) ([]*M, error) {
	var objs []*M
	r.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&objs)
	return objs, nil
}

func (r *PgRepo[M]) GetAll(ctx context.Context) ([]*M, error) {
	var objs []*M
	r.DB.WithContext(ctx).Find(&objs)
	return objs, nil
}

func (r *PgRepo[M]) Create(ctx context.Context, exp *M) (*M, error) {
	if result := r.DB.WithContext(ctx).Create(exp); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *PgRepo[M]) Delete(ctx context.Context, id uint) (*M, error) {
	obj, err := r.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	if result := r.DB.WithContext(ctx).Delete(&obj, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *PgRepo[M]) Update(ctx context.Context, exp *M, values map[string]interface{}) (*M, error) {
	if result := r.DB.WithContext(ctx).Model(&exp).Updates(values); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *PgRepo[M]) CreateMulti(ctx context.Context, exps []*M, batchSize int) ([]*M, error) {
	result := r.DB.WithContext(ctx).Session(&gorm.Session{CreateBatchSize: batchSize}).Create(exps)

	if result.Error != nil {
		return nil, result.Error
	}
	return exps, nil
}

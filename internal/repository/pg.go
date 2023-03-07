package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PgRepo[M any] struct {
	db *gorm.DB
}

func CreatePgRepo[M any](db *gorm.DB) PgRepo[M] {
	return PgRepo[M]{db: db}
}

func (r *PgRepo[M]) Get(ctx context.Context, id uuid.UUID) (res *M, err error) {
	var obj *M
	if result := r.db.WithContext(ctx).First(&obj, "id = ?", id.String()); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *PgRepo[M]) GetMulti(ctx context.Context, limit, offset int) (res []*M, err error) {
	var objs []*M
	r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&objs)
	return objs, nil
}

func (r *PgRepo[M]) Create(ctx context.Context, exp *M) (res *M, err error) {
	if result := r.db.WithContext(ctx).Create(exp); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *PgRepo[M]) Delete(ctx context.Context, id uuid.UUID) (res *M, err error) {
	obj, err := r.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	if result := r.db.WithContext(ctx).Delete(&obj, "id = ?", id.String()); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

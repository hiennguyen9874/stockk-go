package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseRepo[M any] struct{}

func (r *BaseRepo[M]) Get(db *gorm.DB, id uuid.UUID) (res *M, err error) {
	var obj *M
	db.First(&obj, "id = ?", id.String())
	return obj, nil
}

func (r *BaseRepo[M]) GetMulti(db *gorm.DB, limit, offset int) (res []*M, err error) {
	var objs []*M
	db.Limit(limit).Offset(offset).Find(&objs)
	return objs, nil
}

func (r *BaseRepo[M]) Create(db *gorm.DB, exp *M) (res *M, err error) {
	result := db.Create(exp)
	if result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

package gorm

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func List[T schema.Tabler](ctx context.Context, db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) ([]T, error) {
	db = db.WithContext(ctx)
	var list []T
	db = db.Scopes(scopes...).FindInBatches(&list, 100, DefaultFindInBatchesCallback)
	if db.Error != nil {
		return nil, db.Error
	}
	return list, nil
}

func Query[T schema.Tabler](ctx context.Context, db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (T, error) {
	db = db.WithContext(ctx)
	var model, empty T
	db = db.Scopes(scopes...).Find(&model)
	if db.Error != nil {
		return empty, db.Error
	}
	if db.RowsAffected == 0 {
		return empty, nil
	}
	return model, nil
}

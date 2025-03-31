package repository

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Laundry) CreateSetting(ctx context.Context, mSetting *model.Setting) error {
	tx := r.db.WithContext(ctx)

	if err := tx.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{"value": gorm.Expr("VALUES(value)")}),
	}).Create(mSetting).Error; err != nil {
		return err
	}

	return nil
}

func (r *Laundry) UpdateSettingByKey(ctx context.Context, key string, mSetting *model.Setting) error {
	tx := r.db.WithContext(ctx)

	if err := tx.Where(model.Setting{Key: key}).Updates(mSetting).Error; err != nil {
		return err
	}

	return nil
}

func (r *Laundry) ListSetting(ctx context.Context) ([]model.Setting, error) {
	tx := r.db.WithContext(ctx)

	var mSettings []model.Setting
	if err := tx.Find(&mSettings).Error; err != nil {
		return nil, err
	}

	return mSettings, nil
}

func (r *Laundry) DeleteSettingByKey(ctx context.Context, key string) error {
	tx := r.db.WithContext(ctx)

	if err := tx.Where("`key` = ?", key).Delete(&model.Setting{}).Error; err != nil {
		return err
	}

	return nil
}

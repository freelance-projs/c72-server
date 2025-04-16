package repository

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/model"
)

func (r *Repository) UpdateSetting(ctx context.Context, mSetting *model.Setting) error {
	tx := r.db.WithContext(ctx)

	if err := tx.Where("id = ?", 1).Updates(mSetting).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetSetting(ctx context.Context) (*model.Setting, error) {
	tx := r.db.WithContext(ctx)

	var mSetting model.Setting
	if err := tx.Take(&mSetting).Error; err != nil {
		return nil, err
	}

	return &mSetting, nil
}

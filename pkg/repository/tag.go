package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/qb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (t *Repository) UpdateTagNameByName(ctx context.Context, oldName, newName string) error {
	tx := t.db.WithContext(ctx)

	return tx.Transaction(func(tx *gorm.DB) error {
		query := "UPDATE tag SET name = ? WHERE (name = ? OR id = ?)"
		if err := tx.Exec(query, newName, oldName, oldName).Error; err != nil {
			return err
		}
		return nil
	})
}

func (t *Repository) GetTagByID(ctx context.Context, id string) (*model.Tag, error) {
	tx := t.db.WithContext(ctx)

	var mTag model.Tag
	if err := tx.First(&mTag, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrDataNotFound
		}
		return nil, err
	}

	return &mTag, nil
}

func (t *Repository) GetTagsByIDs(ctx context.Context, ids []string) ([]model.Tag, error) {
	tx := t.db.WithContext(ctx)

	var mTags []model.Tag
	if err := tx.Find(&mTags, "id IN ?", ids).Error; err != nil {
		return nil, err
	}

	return mTags, nil
}

func (t *Repository) GetTagsByLendingID(ctx context.Context, lendingID int) ([]model.LendingTag, error) {
	tx := t.db.WithContext(ctx)

	var mLendingTags []model.LendingTag
	if err := tx.Preload(model.LendingTag{}.TagRelation()).
		Where("lending_id = ?", lendingID).
		Find(&mLendingTags).Error; err != nil {
		return nil, err
	}

	return mLendingTags, nil
}

func (t *Repository) GetTagsByWashingID(ctx context.Context, lendingID int) ([]model.LaundryTag, error) {
	tx := t.db.WithContext(ctx)

	var mWashingTags []model.LaundryTag
	if err := tx.Preload(model.LaundryTag{}.TagRelation()).
		Where("laundry_id = ?", lendingID).
		Find(&mWashingTags).Error; err != nil {
		return nil, err
	}

	return mWashingTags, nil
}

func (t *Repository) UpdateTagNameByID(ctx context.Context, mTag *model.Tag) error {
	tx := t.db.WithContext(ctx)

	if err := tx.Updates(mTag).Error; err != nil {
		return err
	}

	return nil
}

func (t *Repository) DeleteTagByID(ctx context.Context, id string) error {
	tx := t.db.WithContext(ctx)

	if err := tx.Exec("DELETE FROM tag WHERE id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateTagInBatches(ctx context.Context, tagIDs []string, name string) error {
	tx := r.db.WithContext(ctx)

	mTags := lodash.Map(tagIDs, func(tagID string, _ int) model.Tag {
		return model.Tag{
			ID:   tagID,
			Name: sql.NullString{String: name, Valid: true},
		}
	})

	return tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{"name": name}),
		}).CreateInBatches(mTags, 100).Error; err != nil {
			return err
		}

		// insert tag_name if not exists
		if err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&model.TagName{Name: name}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *Repository) ListTags(ctx context.Context, filter qb.Builder) ([]model.Tag, error) {
	var results []model.Tag

	tx := r.db.WithContext(ctx)
	if filter != nil {
		tx = filter.Build(tx)
	}
	if err := tx.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

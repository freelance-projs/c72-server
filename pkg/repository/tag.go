package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/apperror"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Tag struct {
	db *gorm.DB
}

func NewTag(db *gorm.DB) *Tag {
	return &Tag{
		db: db,
	}
}

func (t *Tag) UpdateTagName(ctx context.Context, mTags []model.Tag) error {
	tx := t.db.WithContext(ctx)

	query := "UPDATE tags SET name = ? WHERE id = ?"
	txErr := tx.Transaction(func(tx *gorm.DB) error {
		for _, mTag := range mTags {
			if err := tx.Exec(query, mTag.Name, mTag.ID).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return txErr
}

func (t *Tag) ScanTagInBatches(ctx context.Context, mTags []model.Tag) error {
	tx := t.db.WithContext(ctx)

	mTagScanHistories := make([]model.TagScanHistory, 0, len(mTags))
	for _, v := range mTags {
		mTagScanHistories = append(mTagScanHistories, model.TagScanHistory{
			TagID: v.ID,
		})
	}

	txErr := tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(mTags, 100).Error; err != nil {
			return err
		}
		if err := tx.CreateInBatches(mTagScanHistories, 100).Error; err != nil {
			return err
		}

		return nil
	})

	return txErr
}

func (t *Tag) GetTagsScanHistories(ctx context.Context, from, to *time.Time) ([]model.TagScanHistory, error) {
	tx := t.db.WithContext(ctx)

	query := "SELECT * FROM tag_scan_histories WHERE created_at >= ? AND created_at <= ?"
	tx = tx.Raw(query, from, to)

	var mTagScanHistories []model.TagScanHistory
	if err := tx.Scan(&mTagScanHistories).Error; err != nil {
		return nil, err
	}

	return mTagScanHistories, nil
}

func (t *Tag) GetTagByID(ctx context.Context, id string) (*model.Tag, error) {
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

func (t *Tag) GetTagsByIDs(ctx context.Context, ids []string) ([]model.Tag, error) {
	tx := t.db.WithContext(ctx)

	var mTags []model.Tag
	if err := tx.Find(&mTags, "id IN ?", ids).Error; err != nil {
		return nil, err
	}

	return mTags, nil
}

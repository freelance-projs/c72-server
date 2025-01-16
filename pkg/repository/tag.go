package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/qb"
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

func (t *Tag) UpdateTagNameByID(ctx context.Context, mTags []model.Tag) error {
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

func (t *Tag) UpdateTagNameByName(ctx context.Context, oldName, newName string) error {
	tx := t.db.WithContext(ctx)

	query := "UPDATE tags SET name = ? WHERE (name = ? OR id = ?)"
	if err := tx.Exec(query, newName, oldName, oldName).Error; err != nil {
		return err
	}

	return nil
}
func (t *Tag) CreateTagInBatches(ctx context.Context, mTags []model.Tag) error {
	tx := t.db.WithContext(ctx)

	if err := tx.
		Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{"name": gorm.Expr("VALUES(name)")}),
		}).
		CreateInBatches(mTags, 100).Error; err != nil {
		return err
	}

	return nil
}

func (t *Tag) ScanTagInBatches(ctx context.Context, mTags []model.Tag) error {
	tx := t.db.WithContext(ctx)

	// mTagScanHistories := make([]model.TagScanHistory, 0, len(mTags))
	// for _, v := range mTags {
	// 	mTagScanHistories = append(mTagScanHistories, model.TagScanHistory{
	// 		TagID: v.ID,
	// 	})
	// }

	// txErr := tx.Transaction(func(tx *gorm.DB) error {
	if err := tx.
		Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{"is_scanned": true}),
		}).
		CreateInBatches(mTags, 100).Error; err != nil {
		return err
	}
	// if err := tx.CreateInBatches(mTagScanHistories, 100).Error; err != nil {
	// 	return err
	// }
	//
	// 	return nil
	// })

	return nil
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

func (t *Tag) ListTags(ctx context.Context, filter qb.Builder) ([]model.Tag, error) {
	var results []model.Tag

	tx := t.db.WithContext(ctx)
	if filter != nil {
		tx = filter.Build(tx)
	}
	if err := tx.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (t *Tag) GetTags(ctx context.Context, from, to *time.Time) ([]model.Tag, error) {
	tx := t.db.WithContext(ctx)

	query := "SELECT * FROM tags WHERE created_at >= ? AND created_at <= ?"
	tx = tx.Raw(query, from, to)

	var mTags []model.Tag
	if err := tx.Scan(&mTags).Error; err != nil {
		return nil, err
	}

	return mTags, nil
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

func (t *Tag) GetTagsByFilter(ctx context.Context, name string) ([]model.Tag, error) {
	tx := t.db.WithContext(ctx)

	var mTags []model.Tag
	if err := tx.Find(&mTags, "name = ? OR id = ?", name, name).Error; err != nil {
		return nil, err
	}

	return mTags, nil
}

func (t *Tag) DeleteTagByID(ctx context.Context, id string) error {
	tx := t.db.WithContext(ctx)

	if err := tx.Exec("DELETE FROM tags WHERE id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (t *Tag) DeleteTagByName(ctx context.Context, name string) error {
	tx := t.db.WithContext(ctx)

	if err := tx.Exec("DELETE FROM tags WHERE name = ? OR id = ?", name, name).Error; err != nil {
		return err
	}

	return nil
}

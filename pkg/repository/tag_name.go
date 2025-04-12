package repository

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (t *Repository) CreateTagNameInBatches(ctx context.Context, tagNames []model.TagName) error {
	tx := t.db.WithContext(ctx)

	if err := tx.Clauses(clause.OnConflict{
		DoNothing: true,
	}).CreateInBatches(tagNames, 100).Error; err != nil {
		return err
	}

	return nil
}

func (t *Repository) DeleteTagNames(ctx context.Context, names []string) error {
	tx := t.db.WithContext(ctx)

	return tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("name IN ?", names).Delete(&model.TagName{}).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE tag SET name = NULL WHERE name IN ?", names).Error; err != nil {
			return err
		}

		return nil
	})
}

func (t *Repository) UpdateTagName(ctx context.Context, oldName, newName string) error {
	tx := t.db.WithContext(ctx)

	return tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE tag_name SET name = ? WHERE name = ?",
			newName, oldName).Error; err != nil {
			return err
		}

		if err := tx.Exec("UPDATE tag SET name = ? WHERE name = ?",
			newName, oldName).Error; err != nil {
			return err
		}

		return nil
	})
}

func (t *Repository) ListTagNames(ctx context.Context) ([]model.TagName, error) {
	tx := t.db.WithContext(ctx)

	var tagNames []model.TagName
	if err := tx.Find(&tagNames).Error; err != nil {
		return nil, err
	}

	return tagNames, nil
}

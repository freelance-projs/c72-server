package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/qb"
	"gorm.io/gorm"
)

func (r *Laundry) CreateLending(ctx context.Context, mLending *model.Lending, tagIDs []string) error {
	tx := r.db.WithContext(ctx)

	now := time.Now()
	txErr := tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(mLending).Error; err != nil {
			return err
		}

		mLendingTags := lodash.Map(tagIDs, func(id string, i int) model.LendingTag {
			return model.LendingTag{
				LendingID: mLending.ID,
				TagID:     id,
				Status:    model.LendingStatusLending,
			}
		})
		if err := tx.CreateInBatches(mLendingTags, 100).Error; err != nil {
			return err
		}

		// update last_used
		if err := tx.Exec("UPDATE tag SET last_used = ? WHERE id IN (?)", now, tagIDs).Error; err != nil {
			return err
		}

		return nil
	})

	return txErr
}

func (r *Laundry) GetLendingByID(ctx context.Context, id int) (*model.Lending, error) {
	tx := r.db.WithContext(ctx)

	var mLending model.Lending
	if err := tx.Where(model.Lending{ID: id}).
		Preload(mLending.TagsRelation()).
		Take(&mLending).Error; err != nil {
		return nil, err
	}

	return &mLending, nil
}

func (r *Laundry) ReturnDirty(ctx context.Context, tagIDs []string) error {
	tx := r.db.WithContext(ctx)

	txErr := tx.Transaction(func(tx *gorm.DB) error {
		// select lending id
		lendingTagCol := model.LendingTag{}.Columns()
		var mLendingTagIDs []int
		if err := tx.
			Model(model.LendingTag{}).
			Where(fmt.Sprintf("%s IN (?) AND status = ?", lendingTagCol.TagID), tagIDs, model.LendingStatusLending).
			Pluck(lendingTagCol.LendingID, &mLendingTagIDs).Error; err != nil {
			return err
		}

		// map lending id to num tag
		coundLendingID := make(map[int]int)
		for _, id := range mLendingTagIDs {
			coundLendingID[id]++
		}

		// decrement num tag
		for id, numReturned := range coundLendingID {
			if err := tx.Exec("UPDATE lending SET num_lending = num_lending - ? WHERE id = ?", numReturned, id).Error; err != nil {
				return err
			}
		}

		if err := tx.Exec("UPDATE lending_tag SET status = ? WHERE tag_id IN (?)",
			model.LendingStatusReturned, tagIDs).Error; err != nil {
			return err
		}

		return nil
	})

	return txErr
}

func (r *Laundry) ListLending(ctx context.Context, filter qb.Builder) ([]model.Lending, error) {
	tx := r.db.WithContext(ctx)
	if filter != nil {
		tx = filter.Build(tx)
	}

	var mLendings []model.Lending
	if err := tx.Find(&mLendings).Error; err != nil {
		return nil, err
	}

	return mLendings, nil
}

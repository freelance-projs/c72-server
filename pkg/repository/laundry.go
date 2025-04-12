package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/qb"
	"gorm.io/gorm"
)

func (r *Repository) CreateLaundry(ctx context.Context, mLaundry *model.Laundry, tagIDs []string) error {
	tx := r.db.WithContext(ctx)

	now := time.Now()
	txErr := tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(mLaundry).Error; err != nil {
			return err
		}

		mLaundryTags := lodash.Map(tagIDs, func(id string, i int) model.LaundryTag {
			return model.LaundryTag{
				LaundryID: mLaundry.ID,
				TagID:     id,
				Status:    model.LaundryStatusWashing,
			}
		})
		if err := tx.CreateInBatches(mLaundryTags, 100).Error; err != nil {
			return err
		}

		// update last_washing
		if err := tx.Exec("UPDATE tag SET last_washing = ? WHERE id IN (?)", now, tagIDs).Error; err != nil {
			return err
		}

		return nil
	})

	return txErr
}

func (r *Repository) GetLaundryByID(ctx context.Context, id int) (*model.Laundry, error) {
	tx := r.db.WithContext(ctx)

	var mLaundry model.Laundry
	if err := tx.Where(model.Lending{ID: id}).Preload("Tags").Take(&mLaundry).Error; err != nil {
		return nil, err
	}

	return &mLaundry, nil
}

func (r *Repository) ReturnClean(ctx context.Context, tagIDs []string) error {
	tx := r.db.WithContext(ctx)

	txErr := tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE laundry_tag SET status = ? WHERE tag_id IN (?)",
			model.LaundryStatusReturned, tagIDs).Error; err != nil {
			return err
		}

		// select laundry id
		laundryTagCol := model.LaundryTag{}.Columns()
		var mLaundryTagIDs []int
		if err := tx.
			Model(model.LaundryTag{}).
			Where(fmt.Sprintf("%s IN (?)", laundryTagCol.TagID), tagIDs).
			Pluck(laundryTagCol.LaundryID, &mLaundryTagIDs).Error; err != nil {
			return err
		}

		// map laundry id to num tag
		countLaundryID := make(map[int]int)
		for _, id := range mLaundryTagIDs {
			countLaundryID[id]++
		}

		// decrement num tag
		for id, numReturned := range countLaundryID {
			if err := tx.Exec("UPDATE laundry SET num_washing = num_washing - ? WHERE id = ?", numReturned, id).Error; err != nil {
				// Detect this error Error 1690 (22003): BIGINT UNSIGNED value is out of range in
				if apperror.IsMySQLOutOfRange(err) {
					return apperror.New("Không thể trả nhiều đồ giặt hơn số lượng đã giặt")
				}

				return err
			}
		}

		return nil
	})

	return txErr
}

func (r *Repository) ListLaundries(ctx context.Context, filter qb.Builder) ([]model.Laundry, error) {
	tx := r.db.WithContext(ctx)

	if filter != nil {
		tx = filter.Build(tx)
	}

	var laundries []model.Laundry
	if err := tx.Find(&laundries).Error; err != nil {
		return nil, err
	}

	return laundries, nil
}

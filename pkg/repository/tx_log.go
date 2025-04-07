package repository

import (
	"cmp"
	"context"
	"fmt"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/lodash"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Laundry) GetActiveTags(ctx context.Context, action string, tagIDs []string) ([]model.Tag, error) {
	tx := r.db.WithContext(ctx)

	switch action {
	// tag_ids must be not exist in tx_tag
	case "lending":
		var txTagIDs []string
		if err := tx.Select("tag_id").
			Model(model.TxTag{}).
			Where("tag_id IN ?", tagIDs).
			Scan(&txTagIDs).
			Error; err != nil {
			return nil, err
		}
		mapTxTagIDs := make(map[string]struct{}, len(txTagIDs))
		for _, txTagID := range txTagIDs {
			mapTxTagIDs[txTagID] = struct{}{}
		}

		tagIDs = lodash.Filter(tagIDs, func(tagID string, _ int) bool {
			_, exist := mapTxTagIDs[tagID]
			return !exist
		})

		var mTags []model.Tag
		if err := tx.Where("id IN ?", tagIDs).Find(&mTags).Error; err != nil {
			return nil, err
		}
		return mTags, nil

	// tag_ids must be exist in tx_tag with status = lending
	case "lending_return":
		var txTagIDs []string
		if err := tx.Select("tag_id").
			Model(model.TxTag{}).
			Where("tag_id IN ? AND status = ?", tagIDs, model.TxTagStatusLending).
			Scan(&txTagIDs).
			Error; err != nil {
			return nil, err
		}
		var mTags []model.Tag
		if err := tx.Where("id IN ?", txTagIDs).Find(&mTags).Error; err != nil {
			return nil, err
		}
		return mTags, nil

	// tag_ids must be not exist in tx_tag
	case "washing":
		var txTagIDs []string
		if err := tx.Select("tag_id").
			Model(model.TxTag{}).
			Where("tag_id IN ?", tagIDs).
			Scan(&txTagIDs).
			Error; err != nil {
			return nil, err
		}
		mapTxTagIDs := make(map[string]struct{}, len(txTagIDs))
		for _, txTagID := range txTagIDs {
			mapTxTagIDs[txTagID] = struct{}{}
		}

		tagIDs = lodash.Filter(tagIDs, func(tagID string, _ int) bool {
			_, exist := mapTxTagIDs[tagID]
			return !exist
		})

		var mTags []model.Tag
		if err := tx.Where("id IN ?", tagIDs).Find(&mTags).Error; err != nil {
			return nil, err
		}
		return mTags, nil

	// tag_ids must be exist in tx_tag with status = washing
	case "washing_return":
		var txTagIDs []string
		if err := tx.Select("tag_id").
			Model(model.TxTag{}).
			Where("tag_id IN ? AND status = ?", tagIDs, model.TxTagStatusWashing).
			Scan(&txTagIDs).
			Error; err != nil {
			return nil, err
		}
		var mTags []model.Tag
		if err := tx.Where("id IN ?", txTagIDs).Find(&mTags).Error; err != nil {
			return nil, err
		}
		return mTags, nil

	default:
		return nil, apperror.New("action not found")
	}
}

func (r *Laundry) CreateLendingTx(ctx context.Context, department string, tagIDs []string) (*model.TxLogDepartment, error) {
	tx := r.db.WithContext(ctx)

	var mTxLog *model.TxLogDepartment
	txErr := tx.Transaction(func(tx *gorm.DB) error {
		// check tag is free
		var countTxTag int64
		if err := tx.Raw("select count(*) from tx_tag where tag_id IN ? for update", tagIDs).
			Count(&countTxTag).Error; err != nil {
			return err
		}

		if countTxTag > 0 {
			return apperror.New(fmt.Sprintf("Có %d tags không hợp lệ (do đang được mượn hoặc đang được giặt)", countTxTag))
		}

		// get tag name by id
		tags, err := getTagByIDs(tx, tagIDs)
		if err != nil {
			return err
		}
		if len(tags) == 0 {
			return fmt.Errorf("%w (tag)", apperror.DataNotFound)
		}

		// count tag_name
		mapTagName := make(map[string]int)
		for _, tag := range tags {
			mapTagName[cmp.Or(tag.Name.String, tag.ID)]++
		}

		// create tx_log
		txLogTracking := make([]model.TxLogTracking, 0, len(mapTagName))
		for tagName, count := range mapTagName {
			txLogTracking = append(txLogTracking, model.TxLogTracking{
				Name:  tagName,
				Count: count,
			})
		}

		txLog := model.TxLogDepartment{
			Overview: model.TxLogOverview{
				Actor:     department,
				TotalTags: uint(len(tagIDs)),
				Returned:  0,
			},
			Details: []model.TxLogDetail{
				{
					Action:    fmt.Sprintf("%s mượn đồ", department),
					Tracking:  txLogTracking,
					CreatedAt: time.Now(),
				},
			},
		}

		if err := tx.Omit("updated_at").Create(&txLog).Error; err != nil {
			return err
		}
		mTxLog = &txLog

		// create tx_tag to tracking updated tag
		txTags := make([]model.TxTag, 0, len(tagIDs))
		for _, tagID := range tagIDs {
			txTags = append(txTags, model.TxTag{
				TagID:  tagID,
				TxID:   txLog.ID,
				Status: model.TxTagStatusLending,
			})
		}
		if err := tx.Create(&txTags).Error; err != nil {
			return err
		}

		return nil
	})

	return mTxLog, txErr
}

func (r *Laundry) ReturnLendingTx(ctx context.Context, department string, tagIDs []string) ([]int, error) {
	tx := r.db.WithContext(ctx)

	var txIDs []int
	txErr := tx.Transaction(func(tx *gorm.DB) error {
		// check tag is lending
		var mTxTags []model.TxTag
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("tag_id IN ? AND status = ?", tagIDs, model.LendingStatusLending).
			Find(&mTxTags).
			Error; err != nil {
			return err
		}
		if len(mTxTags) != len(tagIDs) {
			return apperror.New(fmt.Sprintf("Yêu cầu tags phải được cho khoa mượn thì mới có thể đem trả, có %d tags không hợp lệ", len(tagIDs)-len(mTxTags)))
		}

		mapTxIDToTagIDs := make(map[int][]string)
		for _, txTag := range mTxTags {
			mapTxIDToTagIDs[txTag.TxID] = append(mapTxIDToTagIDs[txTag.TxID], txTag.TagID)
		}

		now := time.Now()
		for txID, tagIDs := range mapTxIDToTagIDs {
			// get tag name by id
			tags, err := getTagByIDs(tx, tagIDs)
			if err != nil {
				return err
			}
			if len(tags) == 0 {
				return fmt.Errorf("%w (tag)", apperror.DataNotFound)
			}
			// count tag_name
			mapTagName := make(map[string]int)
			for _, tag := range tags {
				mapTagName[cmp.Or(tag.Name.String, tag.ID)]++
			}

			// add tx_log
			txLogTrackings := make([]model.TxLogTracking, 0, len(mapTagName))
			for tagName, count := range mapTagName {
				txLogTrackings = append(txLogTrackings, model.TxLogTracking{
					Name:  tagName,
					Count: count,
				})
			}

			detail := model.TxLogDetail{
				Action:    fmt.Sprintf("%s trả đồ", department),
				Tracking:  txLogTrackings,
				CreatedAt: now,
			}

			txIDs = append(txIDs, txID)
			// get tx_log
			txLog := model.TxLogDepartment{
				ID: txID,
			}
			if err := tx.Where(&txLog).Take(&txLog).Error; err != nil {
				return fmt.Errorf("failed to get tx_log: %w", err)
			}
			// append detail to tx_log
			txLog.Details = append(txLog.Details, detail)
			txLog.Overview.Returned += uint(len(tagIDs))

			// update tx_log
			if err := tx.Updates(&txLog).Error; err != nil {
				return fmt.Errorf("failed to update tx_log: %w", err)
			}

			// clear tx_tag
			if err := tx.Delete(model.TxTag{}, "tag_id IN ?", tagIDs).Error; err != nil {
				return fmt.Errorf("failed to delete tx_tag: %w", err)
			}
		}

		return nil
	})

	return txIDs, txErr
}

func (r *Laundry) CreateWashingTx(ctx context.Context, company string, tagIDs []string) (*model.TxLogCompany, error) {
	tx := r.db.WithContext(ctx)

	var mTxLog *model.TxLogCompany
	txErr := tx.Transaction(func(tx *gorm.DB) error {
		// check tag is free
		var countTxTag int64
		if err := tx.Raw("select count(*) from tx_tag where tag_id IN ? for update", tagIDs).
			Count(&countTxTag).Error; err != nil {
			return err
		}

		if countTxTag > 0 {
			return apperror.New(fmt.Sprintf("Có %d tags không hợp lệ (do đang được mượn hoặc đang được giặt)", countTxTag))
		}

		// get tag name by id
		tags, err := getTagByIDs(tx, tagIDs)
		if err != nil {
			return err
		}
		if len(tags) == 0 {
			return fmt.Errorf("%w (tag)", apperror.DataNotFound)
		}

		// count tag_name
		mapTagName := make(map[string]int)
		for _, tag := range tags {
			mapTagName[cmp.Or(tag.Name.String, tag.ID)]++
		}

		// create tx_log
		txLogTracking := make([]model.TxLogTracking, 0, len(mapTagName))
		for tagName, count := range mapTagName {
			txLogTracking = append(txLogTracking, model.TxLogTracking{
				Name:  tagName,
				Count: count,
			})
		}

		txLog := model.TxLogCompany{
			Overview: model.TxLogOverview{
				Actor:     company,
				TotalTags: uint(len(tagIDs)),
				Returned:  0,
			},
			Details: []model.TxLogDetail{
				{
					Action:    fmt.Sprintf("%s giặt đồ", company),
					Tracking:  txLogTracking,
					CreatedAt: time.Now(),
				},
			},
		}

		if err := tx.Omit("updated_at").Create(&txLog).Error; err != nil {
			return err
		}
		mTxLog = &txLog

		// create tx_tag to tracking updated tag
		txTags := make([]model.TxTag, 0, len(tagIDs))
		for _, tagID := range tagIDs {
			txTags = append(txTags, model.TxTag{
				TagID:  tagID,
				TxID:   txLog.ID,
				Status: model.TxTagStatusWashing,
			})
		}
		if err := tx.Create(&txTags).Error; err != nil {
			return err
		}

		return nil
	})

	return mTxLog, txErr
}

func (r *Laundry) ReturnWashingTx(ctx context.Context, company string, tagIDs []string) ([]int, error) {
	tx := r.db.WithContext(ctx)

	var txIDs []int
	txErr := tx.Transaction(func(tx *gorm.DB) error {
		// check tag is lending
		var mTxTags []model.TxTag
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("tag_id IN ? AND status = ?", tagIDs, model.TxTagStatusWashing).
			Find(&mTxTags).
			Error; err != nil {
			return err
		}
		if len(mTxTags) != len(tagIDs) {
			return apperror.New(fmt.Sprintf("Yêu cầu tags phải được giặt thì mới có thể đem trả, có %d tags không hợp lệ", len(tagIDs)-len(mTxTags)))
		}

		mapTxIDToTagIDs := make(map[int][]string)
		for _, txTag := range mTxTags {
			mapTxIDToTagIDs[txTag.TxID] = append(mapTxIDToTagIDs[txTag.TxID], txTag.TagID)
		}

		now := time.Now()
		for txID, tagIDs := range mapTxIDToTagIDs {
			// get tag name by id
			tags, err := getTagByIDs(tx, tagIDs)
			if err != nil {
				return err
			}
			if len(tags) == 0 {
				return fmt.Errorf("%w (tag)", apperror.DataNotFound)
			}
			// count tag_name
			mapTagName := make(map[string]int)
			for _, tag := range tags {
				mapTagName[cmp.Or(tag.Name.String, tag.ID)]++
			}

			// add tx_log
			txLogTrackings := make([]model.TxLogTracking, 0, len(mapTagName))
			for tagName, count := range mapTagName {
				txLogTrackings = append(txLogTrackings, model.TxLogTracking{
					Name:  tagName,
					Count: count,
				})
			}

			detail := model.TxLogDetail{
				Action:    fmt.Sprintf("%s trả đồ sạch", company),
				Tracking:  txLogTrackings,
				CreatedAt: now,
			}

			txIDs = append(txIDs, txID)
			// get tx_log
			txLog := model.TxLogCompany{
				ID: txID,
			}
			if err := tx.Where(&txLog).Take(&txLog).Error; err != nil {
				return fmt.Errorf("failed to get tx_log: %w", err)
			}
			// append detail to tx_log
			txLog.Details = append(txLog.Details, detail)
			txLog.Overview.Returned += uint(len(tagIDs))

			// update tx_log
			if err := tx.Updates(&txLog).Error; err != nil {
				return fmt.Errorf("failed to update tx_log: %w", err)
			}

			// clear tx_tag
			if err := tx.Delete(model.TxTag{}, "tag_id IN ?", tagIDs).Error; err != nil {
				return fmt.Errorf("failed to delete tx_tag: %w", err)
			}
		}

		return nil
	})

	return txIDs, txErr
}

func getTagByIDs(tx *gorm.DB, tagIDs []string) ([]model.Tag, error) {
	var tags []model.Tag
	if err := tx.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

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

func (r *Repository) GetActiveTags(ctx context.Context, action string, tagIDs []string) ([]model.Tag, error) {
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

		activeTagIDs := lodash.Filter(tagIDs, func(tagID string, _ int) bool {
			_, exist := mapTxTagIDs[tagID]
			return !exist
		})

		var mTags []model.Tag
		if err := tx.Where("id IN ?", activeTagIDs).Find(&mTags).Error; err != nil {
			return nil, err
		}
		if len(mTags) != len(activeTagIDs) {
			return nil, apperror.New("Vui lòng gán tên cho tất cả các tag trước khi thực hiện giao dịch")
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

		activeTagIDs := lodash.Filter(tagIDs, func(tagID string, _ int) bool {
			_, exist := mapTxTagIDs[tagID]
			return !exist
		})

		var mTags []model.Tag
		if err := tx.Where("id IN ?", activeTagIDs).Find(&mTags).Error; err != nil {
			return nil, err
		}

		if len(mTags) != len(activeTagIDs) {
			return nil, apperror.New("Vui lòng gán tên cho tất cả các tag trước khi thực hiện giao dịch")
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

func (r *Repository) CreateLendingTx(ctx context.Context, department string, tagIDs []string) (int, []model.TxLogDepartment, error) {
	tx := r.db.WithContext(ctx)

	txID := int(time.Now().Unix())
	var txLogDepartment []model.TxLogDepartment

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
		if len(tags) != len(tagIDs) {
			return apperror.New("Vui lòng gán tên cho tất cả các tag trước khi thực hiện giao dịch")
		}

		// count tag_name
		mapTagName := make(map[string]int)
		for _, tag := range tags {
			mapTagName[cmp.Or(tag.Name.String, tag.ID)]++
		}

		// create tx_log
		departmentStats := make([]model.TxLogDepartment, 0, len(mapTagName))
		lendingStats := make([]model.LendingStat, 0, len(mapTagName))
		for tagName, count := range mapTagName {
			departmentStats = append(departmentStats, model.TxLogDepartment{
				ID:         txID,
				Department: department,
				Action:     model.EDepartmentActionLending,
				TagName:    tagName,
				Lending:    count,
			})
			lendingStats = append(lendingStats, model.LendingStat{
				ID:         txID,
				TagName:    tagName,
				Lending:    count,
				Department: department,
			})
		}

		if err := createTxLogDepartments(tx, departmentStats); err != nil {
			return err
		}
		if err := createLendingStats(tx, lendingStats); err != nil {
			return err
		}

		// create tx_tag to tracking updated tag
		txTags := make([]model.TxTag, 0, len(tagIDs))
		for _, tagID := range tagIDs {
			txTags = append(txTags, model.TxTag{
				TagID:  tagID,
				TxID:   txID,
				Status: model.TxTagStatusLending,
			})
		}
		if err := tx.Create(&txTags).Error; err != nil {
			return err
		}

		txLogDepartment = departmentStats

		return nil
	})

	return txID, txLogDepartment, txErr
}

func (r *Repository) ReturnLendingTx(ctx context.Context, department string, tagIDs []string) ([]int, []model.TxLogDepartment, error) {
	tx := r.db.WithContext(ctx)

	var txLogDepartments []model.TxLogDepartment

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

		// now := time.Now()
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
			departmentStats := make([]model.TxLogDepartment, 0, len(mapTagName))
			lendingStats := make([]model.LendingStat, 0, len(mapTagName))
			for tagName, count := range mapTagName {
				departmentStats = append(departmentStats, model.TxLogDepartment{
					ID:         txID,
					Department: department,
					Action:     model.EDepartmentActionReturned,
					TagName:    tagName,
					Returned:   count,
				})
				lendingStats = append(lendingStats, model.LendingStat{
					ID:       txID,
					TagName:  tagName,
					Returned: count,
				})
			}

			if err := createTxLogDepartments(tx, departmentStats); err != nil {
				return err
			}
			if err := updateLendingStats(tx, lendingStats); err != nil {
				return err
			}

			txIDs = append(txIDs, txID)

			// clear tx_tag
			if err := tx.Delete(model.TxTag{}, "tag_id IN ?", tagIDs).Error; err != nil {
				return fmt.Errorf("failed to delete tx_tag: %w", err)
			}

			txLogDepartments = append(txLogDepartments, departmentStats...)
		}

		return nil
	})

	return txIDs, txLogDepartments, txErr
}

func (r *Repository) CreateWashingTx(ctx context.Context, company string, tagIDs []string) (int, []model.TxLogCompany, error) {
	tx := r.db.WithContext(ctx)

	var txLogCompany []model.TxLogCompany
	txID := int(time.Now().Unix())

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
		companyStat := make([]model.TxLogCompany, 0, len(mapTagName))
		washingStats := make([]model.WashingStat, 0, len(mapTagName))
		for tagName, count := range mapTagName {
			companyStat = append(companyStat, model.TxLogCompany{
				ID:      txID,
				Company: company,
				Action:  model.ECompanyActionWashing,
				TagName: tagName,
				Washing: count,
			})
			washingStats = append(washingStats, model.WashingStat{
				ID:      txID,
				TagName: tagName,
				Washing: count,
				Company: company,
			})
		}

		if err := createTxLogCompany(tx, companyStat); err != nil {
			return err
		}
		if err := createWashingStats(tx, washingStats); err != nil {
			return err
		}

		// create tx_tag to tracking updated tag
		txTags := make([]model.TxTag, 0, len(tagIDs))
		for _, tagID := range tagIDs {
			txTags = append(txTags, model.TxTag{
				TagID:  tagID,
				TxID:   txID,
				Status: model.TxTagStatusWashing,
			})
		}
		if err := tx.Create(&txTags).Error; err != nil {
			return err
		}

		txLogCompany = companyStat

		return nil
	})

	return txID, txLogCompany, txErr
}

func (r *Repository) ReturnWashingTx(ctx context.Context, company string, tagIDs []string) ([]int, []model.TxLogCompany, error) {
	tx := r.db.WithContext(ctx)

	var txLogCompany []model.TxLogCompany

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
			companyStat := make([]model.TxLogCompany, 0, len(mapTagName))
			washingStats := make([]model.WashingStat, 0, len(mapTagName))
			for tagName, count := range mapTagName {
				companyStat = append(companyStat, model.TxLogCompany{
					ID:       txID,
					Company:  company,
					Action:   model.ECompanyActionReturned,
					TagName:  tagName,
					Returned: count,
				})
				washingStats = append(washingStats, model.WashingStat{
					ID:       txID,
					TagName:  tagName,
					Returned: count,
				})
			}
			if err := createTxLogCompany(tx, companyStat); err != nil {
				return err
			}
			if err := updateWashingStats(tx, washingStats); err != nil {
				return err
			}

			txIDs = append(txIDs, txID)
			// get tx_log
			txLog := model.TxLogCompany{
				ID: txID,
			}
			if err := tx.Where(&txLog).Take(&txLog).Error; err != nil {
				return fmt.Errorf("failed to get tx_log: %w", err)
			}
			// clear tx_tag
			if err := tx.Delete(model.TxTag{}, "tag_id IN ?", tagIDs).Error; err != nil {
				return fmt.Errorf("failed to delete tx_tag: %w", err)
			}

			txLogCompany = append(txLogCompany, companyStat...)
		}

		return nil
	})

	return txIDs, txLogCompany, txErr
}

func getTagByIDs(tx *gorm.DB, tagIDs []string) ([]model.Tag, error) {
	var tags []model.Tag
	if err := tx.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

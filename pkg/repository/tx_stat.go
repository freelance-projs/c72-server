package repository

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/model"
	"gorm.io/gorm"
)

func createLendingStats(tx *gorm.DB, tagStats []model.LendingStat) error {
	return tx.CreateInBatches(tagStats, 20).Error
}

func updateLendingStats(tx *gorm.DB, tagStats []model.LendingStat) error {
	q := `update lending_stat set lending = lending - ?, returned = returned + ? where id = ? and tag_name = ?`
	for _, stat := range tagStats {
		if err := tx.Exec(q, stat.Returned, stat.Returned, stat.ID, stat.TagName).Error; err != nil {
			return err
		}
	}

	return nil
}

func createWashingStats(tx *gorm.DB, tagStats []model.WashingStat) error {
	return tx.CreateInBatches(tagStats, 20).Error
}

func updateWashingStats(tx *gorm.DB, tagStats []model.WashingStat) error {
	q := `update washing_stat set washing = washing - ?, returned = returned + ? where id = ? and tag_name = ?`
	for _, stat := range tagStats {
		if err := tx.Exec(q, stat.Returned, stat.Returned, stat.ID, stat.TagName).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) ListDepartmentStats(ctx context.Context) ([]model.LendingStat, error) {
	tx := r.db.WithContext(ctx)

	q := `select department, sum(lending) as lending, sum(returned) as returned from lending_stat group by department`

	var deptStat []model.LendingStat
	if err := tx.Raw(q).Scan(&deptStat).Error; err != nil {
		return nil, err
	}

	return deptStat, nil
}

func (r *Repository) GetDepartmentStat(ctx context.Context, department string) ([]model.LendingStat, error) {
	tx := r.db.WithContext(ctx)

	q := `select tag_name, sum(lending) as lending , sum(returned) as returned from lending_stat where department = ? group by tag_name`

	var deptStat []model.LendingStat
	if err := tx.Raw(q, department).Scan(&deptStat).Error; err != nil {
		return nil, err
	}

	return deptStat, nil
}

func (r *Repository) GetCompanyStat(ctx context.Context, company string) ([]model.WashingStat, error) {
	tx := r.db.WithContext(ctx)

	q := `select tag_name, sum(washing) as washing , sum(returned) as returned from washing_stat where company = ? group by tag_name`

	var companyStat []model.WashingStat
	if err := tx.Raw(q, company).Scan(&companyStat).Error; err != nil {
		return nil, err
	}

	return companyStat, nil
}

func (r *Repository) ListCompanyStats(ctx context.Context) ([]model.WashingStat, error) {
	tx := r.db.WithContext(ctx)

	q := `select company, sum(washing) as washing, sum(returned) as returned from washing_stat group by company`

	var companyStat []model.WashingStat
	if err := tx.Raw(q).Scan(&companyStat).Error; err != nil {
		return nil, err
	}

	return companyStat, nil
}

func (r *Repository) ListLendingTagStats(ctx context.Context) ([]model.LendingStat, error) {
	tx := r.db.WithContext(ctx)

	q := `select tag_name, sum(lending) as lending, sum(returned) as returned from lending_stat group by tag_name`

	var deptStat []model.LendingStat
	if err := tx.Raw(q).Scan(&deptStat).Error; err != nil {
		return nil, err
	}

	return deptStat, nil
}

func (r *Repository) ListWashingTagStats(ctx context.Context) ([]model.WashingStat, error) {
	tx := r.db.WithContext(ctx)

	q := `select tag_name, sum(washing) as washing, sum(returned) as returned from washing_stat group by tag_name`

	var companyStat []model.WashingStat
	if err := tx.Raw(q).Scan(&companyStat).Error; err != nil {
		return nil, err
	}

	return companyStat, nil
}

func (r *Repository) GetLendingTagStat(ctx context.Context, tagName string) ([]model.LendingStat, error) {
	tx := r.db.WithContext(ctx)

	q := `select department, sum(lending) as lending , sum(returned) as returned from lending_stat where tag_name = ? group by department`

	var tagStat []model.LendingStat
	if err := tx.Raw(q, tagName).Scan(&tagStat).Error; err != nil {
		return nil, err
	}

	return tagStat, nil
}

func (r *Repository) GetWashingTagStat(ctx context.Context, tagName string) ([]model.WashingStat, error) {
	tx := r.db.WithContext(ctx)

	q := `select company, sum(washing) as washing , sum(returned) as returned from washing_stat where tag_name = ? group by company`

	var tagStat []model.WashingStat
	if err := tx.Raw(q, tagName).Scan(&tagStat).Error; err != nil {
		return nil, err
	}

	return tagStat, nil
}

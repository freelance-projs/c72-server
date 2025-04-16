package repository

import (
	"context"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/model"
	"gorm.io/gorm"
)

func (r *Repository) GetTxLogDepartmentByID(ctx context.Context, id int) ([]model.TxLogDepartment, error) {
	tx := r.db.WithContext(ctx)

	var deptLog []model.TxLogDepartment
	if err := tx.Where("id = ?", id).Order("created_at").Find(&deptLog).Error; err != nil {
		return nil, err
	}

	return deptLog, nil
}

func (r *Repository) ListLendStat(ctx context.Context, from time.Time, to time.Time) ([]model.LendingStat, error) {
	tx := r.db.WithContext(ctx)

	q := `select id, department, sum(lending) as lending, sum(returned) as returned, created_at from lending_stat where created_at >= ? and created_at <= ? group by id, department, created_at;`

	var txLogs []model.LendingStat
	if err := tx.Raw(q, from, to).Scan(&txLogs).Error; err != nil {
		return nil, err
	}

	return txLogs, nil
}

func createTxLogDepartments(tx *gorm.DB, departmentStat []model.TxLogDepartment) error {
	return tx.CreateInBatches(departmentStat, 20).Error
}

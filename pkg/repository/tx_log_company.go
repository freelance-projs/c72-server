package repository

import (
	"context"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/model"
	"gorm.io/gorm"
)

func (r *Repository) GetTxLogCompanyByID(ctx context.Context, id int) ([]model.TxLogCompany, error) {
	tx := r.db.WithContext(ctx)

	var txLog []model.TxLogCompany
	if err := tx.Where("id = ?", id).Order("created_at").Find(&txLog).Error; err != nil {
		return nil, err
	}

	return txLog, nil
}

func (r *Repository) ListTxLogCompany(ctx context.Context, from, to time.Time) ([]model.TxLogCompany, error) {
	tx := r.db.WithContext(ctx)

	q := `select id, company, sum(washing) as washing, sum(returned) as returned, created_at from tx_log_company group by id, company, created_at;`

	var txLogs []model.TxLogCompany
	if err := tx.Raw(q).Scan(&txLogs).Error; err != nil {
		return nil, err
	}

	return txLogs, nil
}

func createCompanyStatDetail(tx *gorm.DB, companyStat []model.TxLogCompany) error {
	return tx.Omit("created_at").CreateInBatches(companyStat, 20).Error
}

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

func (r *Repository) ListWashingStat(ctx context.Context, from, to time.Time) ([]model.WashingStat, error) {
	tx := r.db.WithContext(ctx)

	q := `select id, company, sum(washing) as washing, sum(returned) as returned, created_at from washing_stat where created_at >= ? and created_at <= ? group by id, company, created_at;`

	var txLogs []model.WashingStat
	if err := tx.Raw(q, from, to).Scan(&txLogs).Error; err != nil {
		return nil, err
	}

	return txLogs, nil
}

func createTxLogCompany(tx *gorm.DB, companyStat []model.TxLogCompany) error {
	return tx.CreateInBatches(companyStat, 20).Error
}

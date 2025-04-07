package repository

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/qb"
)

func (r *Laundry) ListTxLogDept(ctx context.Context, filter qb.Builder) ([]model.TxLogDepartment, error) {
	tx := r.db.WithContext(ctx)
	if filter != nil {
		tx = filter.Build(tx)
	}

	var txLogs []model.TxLogDepartment
	if err := tx.Omit("details").Find(&txLogs).Error; err != nil {
		return nil, err
	}

	return txLogs, nil
}

func (r *Laundry) GetTxLogDept(ctx context.Context, id int) (*model.TxLogDepartment, error) {
	tx := r.db.WithContext(ctx)

	var txLog model.TxLogDepartment
	if err := tx.Where("id = ?", id).Take(&txLog).Error; err != nil {
		return nil, err
	}

	return &txLog, nil
}

package repository

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/qb"
	"gorm.io/gorm/clause"
)

func (r *Laundry) CreateDepartment(ctx context.Context, mDepartment *model.Department) error {
	tx := r.db.WithContext(ctx)

	if err := tx.Create(mDepartment).Error; err != nil {
		return err
	}

	return nil
}

func (r *Laundry) CreateDepartmentInBatches(ctx context.Context, mDepartments []model.Department) error {
	tx := r.db.WithContext(ctx)

	if err := tx.Clauses(clause.OnConflict{
		DoNothing: true,
	}).CreateInBatches(mDepartments, 100).Error; err != nil {
		return err
	}

	return nil
}

func (r *Laundry) DeleteDepartments(ctx context.Context, names []string) error {
	tx := r.db.WithContext(ctx)

	if err := tx.Where("name IN ?", names).Delete(&model.Department{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *Laundry) ListDepartment(ctx context.Context, filter qb.Builder) ([]model.Department, error) {
	tx := r.db.WithContext(ctx)

	if filter != nil {
		tx = filter.Build(tx)
	}

	var mDepartments []model.Department
	if err := tx.Find(&mDepartments).Error; err != nil {
		return nil, err
	}

	return mDepartments, nil
}

func (r *Laundry) UpdateDepartmentName(ctx context.Context, oldName, newName string) error {
	tx := r.db.WithContext(ctx)

	if err := tx.Exec("UPDATE department SET name = ? WHERE name = ?",
		newName, oldName).Error; err != nil {
		return err
	}

	return nil
}

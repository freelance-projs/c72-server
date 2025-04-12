package repository

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/qb"
	"gorm.io/gorm/clause"
)

func (r *Repository) CreateCompany(ctx context.Context, mCompany *model.Company) error {
	tx := r.db.WithContext(ctx)

	if err := tx.Create(mCompany).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateCompanyInBatches(ctx context.Context, mCompanys []model.Company) error {
	tx := r.db.WithContext(ctx)

	if err := tx.Clauses(clause.OnConflict{
		DoNothing: true,
	}).CreateInBatches(mCompanys, 100).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteCompanies(ctx context.Context, names []string) error {
	tx := r.db.WithContext(ctx)

	if err := tx.Where("name IN ?", names).Delete(&model.Company{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListCompanies(ctx context.Context, filter qb.Builder) ([]model.Company, error) {
	tx := r.db.WithContext(ctx)

	if filter != nil {
		tx = filter.Build(tx)
	}

	var mCompanys []model.Company
	if err := tx.Find(&mCompanys).Error; err != nil {
		return nil, err
	}

	return mCompanys, nil
}

func (r *Repository) UpdateCompanyName(ctx context.Context, oldName, newName string) error {
	tx := r.db.WithContext(ctx)

	if err := tx.Exec("UPDATE company SET name = ? WHERE name = ?",
		newName, oldName).Error; err != nil {
		return err
	}

	return nil
}

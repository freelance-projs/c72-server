package company

import (
	"context"
	"fmt"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
	"github.com/ngoctd314/common/qb"
)

type listCompanyRepo interface {
	ListCompanies(ctx context.Context, filter qb.Builder) ([]model.Company, error)
}

type list struct {
	repo listCompanyRepo
}

func List(repo listCompanyRepo) *list {
	return &list{repo: repo}
}

func (uc *list) Usecase(ctx context.Context, req *dto.ListCompaniesRequest) (*ghttp.ResponseBody, error) {
	mCompanies, err := uc.repo.ListCompanies(ctx, nil)

	if err != nil {
		return nil, fmt.Errorf("list companies error: %w", apperror.GormTranslator(err))
	}

	companies := lodash.Map(mCompanies, func(m model.Company, _ int) dto.Company {
		return mapper.ToCompanyDto(&m)
	})

	return ghttp.ResponseBodyOK(companies), nil
}

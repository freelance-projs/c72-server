package department

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

type listDepartmentRepo interface {
	ListDepartment(ctx context.Context, filter qb.Builder) ([]model.Department, error)
}

type list struct {
	repo listDepartmentRepo
}

func List(repo listDepartmentRepo) *list {
	return &list{repo: repo}
}

func (uc *list) Usecase(ctx context.Context, req *dto.ListDepartmentsRequest) (*ghttp.ResponseBody, error) {
	mDepartments, err := uc.repo.ListDepartment(ctx, nil)

	if err != nil {
		return nil, fmt.Errorf("list departments error: %w", apperror.GormTranslator(err))
	}

	departments := lodash.Map(mDepartments, func(m model.Department, _ int) dto.Department {
		return mapper.ToDepartmentDto(&m)
	})

	return ghttp.ResponseBodyOK(departments), nil
}

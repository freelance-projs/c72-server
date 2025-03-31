package department

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/net/ghttp"
)

type createDepartmentRepo interface {
	CreateDepartment(ctx context.Context, mDepartment *model.Department) error
}

type create struct {
	repo createDepartmentRepo
}

func Create(repo createDepartmentRepo) *create {
	return &create{
		repo: repo,
	}
}

func (uc *create) Usecase(ctx context.Context, req *dto.CreateDepartmentRequest) (*ghttp.ResponseBody, error) {
	mDepartment := &model.Department{Name: req.Name}
	if err := uc.repo.CreateDepartment(ctx, mDepartment); err != nil {
		if apperror.IsMySQLDuplicate(err) {
			return nil, apperror.ErrConflict("department is already exists")
		}
		return nil, err
	}

	return ghttp.ResponseBodyCreated(mapper.ToDepartmentDto(mDepartment), "department"), nil
}

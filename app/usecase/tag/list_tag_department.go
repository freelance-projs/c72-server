package tag

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type listTagDepartment struct {
	repo *repository.Repository
}

func ListTagDepartment(repo *repository.Repository) *listTagDepartment {
	return &listTagDepartment{
		repo: repo,
	}
}

func (uc *listTagDepartment) Usecase(ctx context.Context, req *dto.CreateTagCompanyRequest) (
	*ghttp.ResponseBody, error) {

	mTags, err := uc.repo.ListTagDepartments(ctx)
	if err != nil {
		return nil, err
	}

	dtoTags := lodash.Map(mTags, func(tag model.TagDepartment, _ int) dto.TagDepartment {
		return dto.TagDepartment{
			ID:   tag.ID,
			Name: tag.Name,
		}
	})

	return ghttp.ResponseBodyOK(dtoTags), nil
}

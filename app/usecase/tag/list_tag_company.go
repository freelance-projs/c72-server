package tag

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type listTagCompany struct {
	repo *repository.Repository
}

func ListTagCompanies(repo *repository.Repository) *listTagCompany {
	return &listTagCompany{
		repo: repo,
	}
}

func (uc *listTagCompany) Usecase(ctx context.Context, req *dto.CreateTagCompanyRequest) (
	*ghttp.ResponseBody, error) {

	mTags, err := uc.repo.ListTagCompanies(ctx)
	if err != nil {
		return nil, err
	}

	dtoTags := lodash.Map(mTags, func(tag model.TagCompany, _ int) dto.TagCompany {
		return dto.TagCompany{
			ID:   tag.ID,
			Name: tag.Name,
		}
	})

	return ghttp.ResponseBodyOK(dtoTags), nil
}

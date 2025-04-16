package tag

import (
	"context"
	"fmt"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/net/ghttp"
)

type createTagCompany struct {
	repo *repository.Repository
}

func CreateTagCompany(repo *repository.Repository) *createTagCompany {
	return &createTagCompany{
		repo: repo,
	}
}

func (uc *createTagCompany) Usecase(ctx context.Context, req *dto.CreateTagCompanyRequest) (
	*ghttp.ResponseBody, error) {

	if err := uc.repo.CreateTagCompany(ctx, &model.TagCompany{
		ID:   req.ID,
		Name: req.Name,
	}); err != nil {
		return nil, err
	}

	return ghttp.ResponseBodyCreated(nil, fmt.Sprintf("Gán tag cho %s thành công", req.Name)), nil
}

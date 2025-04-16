package tag

import (
	"context"
	"fmt"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/net/ghttp"
)

type createTagDepartment struct {
	repo *repository.Repository
}

func CreateTagDepartment(repo *repository.Repository) *createTagDepartment {
	return &createTagDepartment{
		repo: repo,
	}
}

func (uc *createTagDepartment) Usecase(ctx context.Context, req *dto.CreateTagCompanyRequest) (
	*ghttp.ResponseBody, error) {

	if err := uc.repo.CreateTagDepartment(ctx, &model.TagDepartment{
		ID:   req.ID,
		Name: req.Name,
	}); err != nil {
		return nil, err
	}

	return ghttp.ResponseBodyCreated(nil, fmt.Sprintf("Gán tag cho %s thành công", req.Name)), nil
}

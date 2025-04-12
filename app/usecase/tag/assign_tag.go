package tag

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/net/ghttp"
)

type assignTag struct {
	repo *repository.Repository
}

func AssignTag(repo *repository.Repository) *assignTag {
	return &assignTag{
		repo: repo,
	}
}

func (uc *assignTag) Usecase(ctx context.Context, req *dto.AssignTagRequest) (*ghttp.ResponseBody, error) {

	if err := uc.repo.CreateTagInBatches(ctx, req.IDs, req.Name); err != nil {
		return nil, err
	}

	return ghttp.ResponseBodyCreated(nil, "Gán tên cho tag thành công"), nil
}

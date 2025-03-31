package lending

import (
	"context"
	"fmt"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/net/ghttp"
)

type returnDirtyRepo interface {
	ReturnDirty(ctx context.Context, tagIDs []string) error
}

type returnDirty struct {
	repo returnDirtyRepo
}

func ReturnDirty(repo returnDirtyRepo) *returnDirty {
	return &returnDirty{
		repo: repo,
	}
}

func (uc *returnDirty) Usecase(ctx context.Context, req *dto.ReturnDirtyRequest) (*ghttp.ResponseBody, error) {
	if err := uc.repo.ReturnDirty(ctx, req.TagIDs); err != nil {
		return nil, fmt.Errorf("return dirty error: %w", apperror.GormTranslator(err))
	}

	return ghttp.ResponseBodyOK("Trả đồ bẩn thành công"), nil
}

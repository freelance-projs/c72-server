package laundry

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/net/ghttp"
)

type returnCleanRepo interface {
	ReturnClean(ctx context.Context, tagIDs []string) error
}

type returnClean struct {
	repo returnCleanRepo
}

func ReturnClean(repo returnCleanRepo) *returnClean {
	return &returnClean{
		repo: repo,
	}
}

func (uc *returnClean) Usecase(ctx context.Context, req *dto.ReturnCleanRequest) (*ghttp.ResponseBody, error) {
	if err := uc.repo.ReturnClean(ctx, req.TagIDs); err != nil {
		return nil, apperror.GormTranslator(err)
	}

	return ghttp.ResponseBodyOK(nil, ghttp.ResponseBodyWithMessage("Trả đồ sạch thành công")), nil
}

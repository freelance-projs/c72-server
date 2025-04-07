package company

import (
	"context"
	"log/slog"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/common/net/ghttp"
)

type updateCompanyRepo interface {
	UpdateCompanyName(ctx context.Context, oldName, newName string) error
}

type update struct {
	repo updateCompanyRepo
}

func Change(repo updateCompanyRepo) *update {
	return &update{
		repo: repo,
	}
}

func (uc *update) Usecase(ctx context.Context, req *dto.ChangeCompanyRequest) (*ghttp.ResponseBody, error) {
	if err := uc.repo.UpdateCompanyName(ctx, req.OldName, req.NewName); err != nil {
		slog.Error("error change tag names", "err", err)
		return nil, err
	}

	return ghttp.ResponseBodyOK("thay đổi tên công ty thành công"), nil
}

package setting

import (
	"context"
	"fmt"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/net/ghttp"
)

type get struct {
	repo *repository.Repository
}

func Get(repo *repository.Repository) *get {
	return &get{repo: repo}
}

func (uc *get) Usecase(ctx context.Context, req *dto.ListTagNameRequest) (*ghttp.ResponseBody, error) {
	mSetting, err := uc.repo.GetSetting(ctx)

	if err != nil {
		return nil, fmt.Errorf("get setting error: %w", apperror.GormTranslator(err))
	}

	return ghttp.ResponseBodyOK(dto.Setting{
		TxLogSheetID:  mSetting.TxLogSheetID,
		ReportSheetID: mSetting.ReportSheetID,
	}), nil
}

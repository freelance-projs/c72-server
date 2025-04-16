package setting

import (
	"context"
	"fmt"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/net/ghttp"
)

type update struct {
	repo *repository.Repository
}

func Update(repo *repository.Repository) *update {
	return &update{repo: repo}
}

func (uc *update) Usecase(ctx context.Context, req *dto.UpdateSettingRequest) (*ghttp.ResponseBody, error) {
	err := uc.repo.UpdateSetting(ctx, &model.Setting{
		TxLogSheetID:  req.TxLogSheetID,
		ReportSheetID: req.ReportSheetID,
	})

	if err != nil {
		return nil, fmt.Errorf("update setting error: %w", apperror.GormTranslator(err))
	}

	return ghttp.ResponseBodyOK(nil), nil
}

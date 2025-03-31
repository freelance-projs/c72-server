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
	repo *repository.Laundry
}

func Update(repo *repository.Laundry) *update {
	return &update{repo: repo}
}

func (uc *update) Usecase(ctx context.Context, req *dto.UpdateSettingRequest) (*ghttp.ResponseBody, error) {
	err := uc.repo.UpdateSettingByKey(ctx, req.Key, &model.Setting{
		Key:   req.Key,
		Value: req.Value,
	})

	if err != nil {
		return nil, fmt.Errorf("update setting error: %w", apperror.GormTranslator(err))
	}

	return ghttp.ResponseBodyOK(nil), nil
}

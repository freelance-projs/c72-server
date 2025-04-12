package setting

import (
	"context"
	"fmt"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type list struct {
	repo *repository.Repository
}

func List(repo *repository.Repository) *list {
	return &list{repo: repo}
}

func (uc *list) Usecase(ctx context.Context, req *dto.ListTagNameRequest) (*ghttp.ResponseBody, error) {
	mSettings, err := uc.repo.ListSetting(ctx)

	if err != nil {
		return nil, fmt.Errorf("list setting error: %w", apperror.GormTranslator(err))
	}

	settingDtos := lodash.Map(mSettings, func(m model.Setting, _ int) dto.Setting {
		return mapper.ToSettingDto(&m)
	})

	return ghttp.ResponseBodyOK(settingDtos), nil
}

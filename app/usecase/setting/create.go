package setting

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/net/ghttp"
)

type create struct {
	repo *repository.Repository
}

func Create(repo *repository.Repository) *create {
	return &create{
		repo: repo,
	}
}

func (uc *create) Usecase(ctx context.Context, req *dto.CreateSettingRequest) (*ghttp.ResponseBody, error) {
	mSetting := model.Setting{
		Key:   req.Key,
		Value: req.Value,
	}

	if err := uc.repo.CreateSetting(ctx, &mSetting); err != nil {
		return nil, err
	}

	settingDto := mapper.ToSettingDto(&mSetting)

	return ghttp.ResponseBodyCreated(settingDto, "setting is created"), nil
}

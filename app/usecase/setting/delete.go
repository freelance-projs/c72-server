package setting

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/net/ghttp"
)

type delete struct {
	repo *repository.Repository
}

func Delete(repo *repository.Repository) *delete {
	return &delete{repo: repo}
}

func (uc *delete) Usecase(ctx context.Context, req *dto.DeleteSettingRequest) (*ghttp.ResponseBody, error) {
	err := uc.repo.DeleteSettingByKey(ctx, req.Key)

	if err != nil {
		return nil, fmt.Errorf("delete setting error: %w", apperror.GormTranslator(err))
	}

	return ghttp.ResponseBodyOK(nil), nil
}

func (uc *delete) Bind(c *gin.Context) (*dto.DeleteSettingRequest, error) {
	var req dto.DeleteSettingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

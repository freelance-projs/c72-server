package company

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/common/net/ghttp"
)

type deleteRepo interface {
	DeleteCompanies(ctx context.Context, tagNames []string) error
}

type deleteBatch struct {
	repo deleteRepo
}

func DeleteBatch(repo deleteRepo) *deleteBatch {
	return &deleteBatch{
		repo: repo,
	}
}

func (uc *deleteBatch) Usecase(ctx context.Context, req *dto.DeleteCompanyRequest) (*ghttp.ResponseBody, error) {
	if err := uc.repo.DeleteCompanies(ctx, req.Names); err != nil {
		slog.Error("error deleting company", "err", err)
		return nil, err
	}

	return ghttp.ResponseBodyOK("companies are deleted"), nil
}

// TODO: check why default delete Bind is not working
func (t *deleteBatch) Bind(c *gin.Context) (*dto.DeleteCompanyRequest, error) {
	var req dto.DeleteCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

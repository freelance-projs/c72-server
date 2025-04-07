package department

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/common/net/ghttp"
)

type deleteBatchTagNameRepo interface {
	DeleteDepartments(ctx context.Context, tagNames []string) error
}

type deleteBatch struct {
	repo deleteBatchTagNameRepo
}

func DeleteBatch(repo deleteBatchTagNameRepo) *deleteBatch {
	return &deleteBatch{
		repo: repo,
	}
}

func (uc *deleteBatch) Usecase(ctx context.Context, req *dto.DeleteDepartmentRequest) (*ghttp.ResponseBody, error) {
	if err := uc.repo.DeleteDepartments(ctx, req.Names); err != nil {
		slog.Error("error deleting department", "err", err)
		return nil, err
	}

	return ghttp.ResponseBodyOK("departments are deleted"), nil
}

// TODO: check why default delete Bind is not working
func (t *deleteBatch) Bind(c *gin.Context) (*dto.DeleteDepartmentRequest, error) {
	var req dto.DeleteDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

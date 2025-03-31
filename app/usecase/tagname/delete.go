package tagname

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/common/net/ghttp"
)

type deleteBatchTagNameRepo interface {
	DeleteTagNames(ctx context.Context, tagNames []string) error
}

type deleteBatch struct {
	repo deleteBatchTagNameRepo
}

func DeleteBatch(repo deleteBatchTagNameRepo) *deleteBatch {
	return &deleteBatch{
		repo: repo,
	}
}

func (uc *deleteBatch) Usecase(ctx context.Context, req *dto.DeleteTagNameRequest) (*ghttp.ResponseBody, error) {
	if err := uc.repo.DeleteTagNames(ctx, req.Names); err != nil {
		slog.Error("error deleting tag names", "err", err)
		return nil, err
	}

	return ghttp.ResponseBodyOK("tag names are deleted"), nil
}

// TODO: check why default delete Bind is not working
func (t *deleteBatch) Bind(c *gin.Context) (*dto.DeleteTagNameRequest, error) {
	var req dto.DeleteTagNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

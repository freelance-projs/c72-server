package tagname

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/common/net/ghttp"
)

type updateTagNameRepo interface {
	UpdateTagName(ctx context.Context, oldName, newName string) error
}

type update struct {
	repo updateTagNameRepo
}

func Change(repo updateTagNameRepo) *update {
	return &update{
		repo: repo,
	}
}

func (uc *update) Usecase(ctx context.Context, req *dto.ChangeTagNameRequest) (*ghttp.ResponseBody, error) {
	if err := uc.repo.UpdateTagName(ctx, req.OldName, req.NewName); err != nil {
		slog.Error("error change tag names", "err", err)
		return nil, err
	}

	return ghttp.ResponseBodyOK("thay đổi tag_name thành công"), nil
}

// TODO: check why default delete Bind is not working
func (t *update) Bind(c *gin.Context) (*dto.ChangeTagNameRequest, error) {
	var req dto.ChangeTagNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

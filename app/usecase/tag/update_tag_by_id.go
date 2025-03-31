package tag

import (
	"context"
	"database/sql"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/net/ghttp"
)

type tagRepositoryUpdateTagname interface {
	UpdateTagNameByID(ctx context.Context, mTag *model.Tag) error
}

type updateTagName struct {
	tagRepo tagRepositoryUpdateTagname
}

func UpdateTagName(tagRepo tagRepositoryUpdateTagname) *updateTagName {
	return &updateTagName{tagRepo}
}

func (uc *updateTagName) Usecase(ctx context.Context, req *dto.UpdateTagRequest) (*ghttp.ResponseBody, error) {
	if err := uc.tagRepo.UpdateTagNameByID(ctx, &model.Tag{
		ID:   req.ID,
		Name: sql.NullString{String: req.Name, Valid: true},
	}); err != nil {
		return nil, err
	}

	return ghttp.ResponseBodyOK(nil), nil
}

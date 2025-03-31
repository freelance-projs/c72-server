package tag

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/common/net/ghttp"
)

type tagRepositoryUpdateTagnameByName interface {
	UpdateTagNameByName(ctx context.Context, oldName, newName string) error
}

type updateTagNameByName struct {
	tagRepo tagRepositoryUpdateTagnameByName
}

func UpdateTagNameByName(tagRepo tagRepositoryUpdateTagnameByName) *updateTagNameByName {
	return &updateTagNameByName{
		tagRepo: tagRepo,
	}
}

func (uc *updateTagNameByName) Usecase(ctx context.Context, req *dto.UpdateTagNameRequestByName) (*ghttp.ResponseBody, error) {
	if err := uc.tagRepo.UpdateTagNameByName(ctx, req.OldName, req.NewName); err != nil {
		return nil, err
	}

	return ghttp.ResponseBodyOK(nil), nil
}

package tag

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/common/net/ghttp"
)

type tagRepositoryDeleteByID interface {
	DeleteTagByID(ctx context.Context, id string) error
}

type deleteTagByID struct {
	tagRepo tagRepositoryDeleteByID
}

func DeleteByID(tagRepo tagRepositoryDeleteByID) *deleteTagByID {
	return &deleteTagByID{
		tagRepo: tagRepo,
	}
}

func (uc *deleteTagByID) Usecase(ctx context.Context, req *dto.DeleteTagByIDRequest) (*ghttp.ResponseBody, error) {
	err := uc.tagRepo.DeleteTagByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return ghttp.ResponseBodyOK(nil), nil
}

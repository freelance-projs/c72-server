package tag

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/net/ghttp"
)

type tagRepoGetByID interface {
	GetTagByID(ctx context.Context, id string) (*model.Tag, error)
}

type getTagByID struct {
	tagRepo tagRepoGetByID
}

func GetTagByIDUsecase(tagRepo tagRepoGetByID) *getTagByID {
	return &getTagByID{
		tagRepo: tagRepo,
	}
}

func (uc *getTagByID) Usecase(ctx context.Context, req *dto.GetTagByIDRequest) (*ghttp.ResponseBody, error) {
	mTag, err := uc.tagRepo.GetTagByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return ghttp.ResponseBodyOK(mapper.ToTagDto(mTag)), nil
}

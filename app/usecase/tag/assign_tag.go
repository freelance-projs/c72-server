package tag

import (
	"context"
	"database/sql"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type tagRepoAssignTag interface {
	CreateTagInBatches(ctx context.Context, mTags []model.Tag) error
}

type assignTag struct {
	tagRepo tagRepoAssignTag
}

func AssignTag(tagRepo tagRepoAssignTag) *assignTag {
	return &assignTag{
		tagRepo: tagRepo,
	}
}

func (uc *assignTag) Usecase(ctx context.Context, req *dto.AssignTagRequest) (*ghttp.ResponseBody, error) {
	mTags := lodash.Map(req.IDs, func(id string, i int) model.Tag {
		return model.Tag{
			ID:   id,
			Name: sql.NullString{String: req.Name, Valid: true},
		}
	})

	if err := uc.tagRepo.CreateTagInBatches(ctx, mTags); err != nil {
		return nil, err
	}

	return ghttp.ResponseBodyCreated(nil, "Gán tên cho tag thành công"), nil
}

package laundry

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type getTagsRepository interface {
	GetTagsByWashingID(ctx context.Context, lendingID int) ([]model.LaundryTag, error)
}

type getTags struct {
	repo getTagsRepository
}

func GetTags(repo getTagsRepository) *getTags {
	return &getTags{repo}
}

func (uc *getTags) Usecase(ctx context.Context, req *dto.GetLendingTagsRequest) (*ghttp.ResponseBody, error) {
	mTags, err := uc.repo.GetTagsByWashingID(ctx, req.LendingID)
	if err != nil {
		return nil, err
	}

	mWashingTagDTOs := lodash.Map(mTags, func(mTag model.LaundryTag, _ int) dto.LaundryTag {
		washingTagDto := mapper.ToLaundryTagDto(&mTag)
		if mTag.Tag != nil {
			washingTagDto.TagName = mTag.Tag.Name.String
		}
		return washingTagDto
	})

	return ghttp.ResponseBodyOK(mWashingTagDTOs), nil
}

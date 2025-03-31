package lending

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type getTagsRepository interface {
	GetTagsByLendingID(ctx context.Context, lendingID int) ([]model.LendingTag, error)
}

type getTags struct {
	repo getTagsRepository
}

func GetTags(repo getTagsRepository) *getTags {
	return &getTags{repo}
}

func (uc *getTags) Usecase(ctx context.Context, req *dto.GetLendingTagsRequest) (*ghttp.ResponseBody, error) {
	mTags, err := uc.repo.GetTagsByLendingID(ctx, req.LendingID)
	if err != nil {
		return nil, err
	}

	mlendingTagDTOs := lodash.Map(mTags, func(mTag model.LendingTag, _ int) dto.LendingTag {
		lendingTagDto := mapper.ToLendingTagDTO(&mTag)
		if mTag.Tag != nil {
			lendingTagDto.TagName = mTag.Tag.Name.String
		}
		return lendingTagDto
	})

	return ghttp.ResponseBodyOK(mlendingTagDTOs), nil
}

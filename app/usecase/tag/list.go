package tag

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/net/ghttp"
	"github.com/ngoctd314/common/qb"
)

type tagRepoFilter interface {
	ListTags(ctx context.Context, filter qb.Builder) ([]model.Tag, error)
}

type listTags struct {
	tagRepo tagRepoFilter
}

func List(tagRepo tagRepoFilter) *listTags {
	return &listTags{
		tagRepo: tagRepo,
	}
}

func (uc *listTags) Usecase(ctx context.Context, req *dto.ListTagsRequest) (*ghttp.ResponseBody, error) {
	filters := uc.filters(req)

	mTags, err := uc.tagRepo.ListTags(ctx, qb.And(filters...))
	if err != nil {
		return nil, err
	}

	tagDTOs := make([]dto.Tag, 0, len(mTags))
	for _, v := range mTags {
		tagDTOs = append(tagDTOs, mapper.ToTagDto(&v))
	}

	return ghttp.ResponseBodyOK(tagDTOs), nil
}

func (uc *listTags) filters(req *dto.ListTagsRequest) []*qb.Cond {
	var filters []*qb.Cond
	if req.Name != nil {
		filters = append(filters, qb.Or(qb.Eq("name", *req.Name), qb.Eq("id", *req.Name)))
	}
	if req.From != nil {
		filters = append(filters, qb.Gte("created_at", *req.From))
	}
	if req.To != nil {
		filters = append(filters, qb.Lte("created_at", *req.To))
	}

	return filters
}

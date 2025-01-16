package tag

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/gvalidator"
	"github.com/ngoctd314/common/qb"
)

type tagRepositoryFilter interface {
	GetTagsByFilter(ctx context.Context, name string) ([]model.Tag, error)
	ListTags(ctx context.Context, filter qb.Builder) ([]model.Tag, error)
}

type listTags struct {
	tagRepo tagRepositoryFilter
}

func List(tagRepo tagRepositoryFilter) gin.HandlerFunc {
	return func(c *gin.Context) {
		uc := &listTags{tagRepo: tagRepo}

		req, bindErr := uc.bind(c)
		if bindErr != nil {
			dto.JSONFail(c, bindErr)
			return
		}

		if validateErr := uc.validate(req); validateErr != nil {
			dto.JSONFail(c, validateErr)
			return
		}

		resp, err := uc.usecase(c.Request.Context(), req)
		if err != nil {
			dto.JSONFail(c, err)
			return
		}

		dto.JSONSuccess(c, resp)
	}
}

func (uc *listTags) usecase(ctx context.Context, req *dto.ListTagsRequest) (*dto.Response, error) {
	filters := uc.filters(req)

	mTags, err := uc.tagRepo.ListTags(ctx, qb.And(filters...))
	if err != nil {
		return nil, err
	}

	tagDTOs := make([]dto.Tag, 0, len(mTags))
	for _, v := range mTags {
		tagDTOs = append(tagDTOs, mapper.ToTagDTO(&v))
	}

	return dto.StatusOK(tagDTOs), nil
}

func (uc *listTags) filters(req *dto.ListTagsRequest) []*qb.Cond {
	var filters []*qb.Cond
	if req.IsScanned != nil {
		filters = append(filters, qb.Eq("is_scanned", req.IsScanned))
	}
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

func (uc *listTags) bind(c *gin.Context) (*dto.ListTagsRequest, error) {
	var req dto.ListTagsRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (uc *listTags) validate(req *dto.ListTagsRequest) error {
	if err := gvalidator.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}

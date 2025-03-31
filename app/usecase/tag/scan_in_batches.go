package tag

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/gvalidator"
)

type tagRepositoryScan interface {
	ScanTagInBatches(ctx context.Context, mTags []model.Tag) error
}

type scanTagInBatches struct {
	tagRepo tagRepositoryScan
}

func ScanInBatches(tagRepo tagRepositoryScan) gin.HandlerFunc {
	return func(c *gin.Context) {
		uc := &scanTagInBatches{
			tagRepo: tagRepo,
		}

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

func (uc *scanTagInBatches) usecase(ctx context.Context, req *dto.ScanTagRequest) (*dto.Response, error) {
	mTags := make([]model.Tag, 0, len(req.TagIDs))
	for _, v := range req.TagIDs {
		mTags = append(mTags, model.Tag{
			ID: v,
		})
	}

	if err := uc.tagRepo.ScanTagInBatches(ctx, mTags); err != nil {
		return nil, err
	}

	tagDTOs := make([]dto.Tag, 0, len(mTags))
	for _, v := range mTags {
		tagDTOs = append(tagDTOs, mapper.ToTagDto(&v))
	}

	return dto.StatusCreated(tagDTOs, "tag scan history"), nil
}

func (uc *scanTagInBatches) bind(c *gin.Context) (*dto.ScanTagRequest, error) {
	var req dto.ScanTagRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (uc *scanTagInBatches) validate(req *dto.ScanTagRequest) error {
	if err := gvalidator.ValidateStruct(req); err != nil {
		return err
	}
	return nil
}

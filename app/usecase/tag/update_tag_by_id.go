package tag

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/gvalidator"
)

type tagRepositoryUpdateTagname interface {
	UpdateTagNameByID(ctx context.Context, mTags []model.Tag) error
}

type updateTagName struct {
	tagRepo tagRepositoryUpdateTagname
}

func UpdateTagName(tagRepo tagRepositoryUpdateTagname) gin.HandlerFunc {
	return func(c *gin.Context) {
		uc := &updateTagName{
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

func (uc *updateTagName) usecase(ctx context.Context, req []dto.UpdateTagNameRequest) (*dto.Response, error) {
	mTags := make([]model.Tag, 0, len(req))
	for _, v := range req {
		mTags = append(mTags, model.Tag{
			ID:   v.ID,
			Name: sql.NullString{String: v.Name, Valid: true},
		})
	}
	if err := uc.tagRepo.UpdateTagNameByID(ctx, mTags); err != nil {
		return nil, err
	}

	return dto.StatusOK(nil), nil
}

func (uc *updateTagName) bind(c *gin.Context) ([]dto.UpdateTagNameRequest, error) {
	req := []dto.UpdateTagNameRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func (uc *updateTagName) validate(req []dto.UpdateTagNameRequest) error {
	if err := gvalidator.ValidateArray(req); err != nil {
		return err
	}

	return nil
}

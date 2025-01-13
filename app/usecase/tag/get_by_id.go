package tag

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/gvalidator"
)

type tagRepositoryGetByID interface {
	GetTagByID(ctx context.Context, id string) (*model.Tag, error)
}

type getTagByID struct {
	tagRepo tagRepositoryGetByID
}

func GetByID(tagRepo tagRepositoryGetByID) gin.HandlerFunc {
	return func(c *gin.Context) {
		uc := &getTagByID{
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

func (uc *getTagByID) usecase(ctx context.Context, req *dto.GetTagByIDRequest) (*dto.Response, error) {
	mTag, err := uc.tagRepo.GetTagByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &dto.Response{
		Success: true,
		Data:    mapper.ToTagDTO(mTag),
	}, nil
}

func (uc *getTagByID) bind(c *gin.Context) (*dto.GetTagByIDRequest, error) {
	var req dto.GetTagByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (uc *getTagByID) validate(req *dto.GetTagByIDRequest) error {
	if err := gvalidator.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}

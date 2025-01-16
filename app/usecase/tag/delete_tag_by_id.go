package tag

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/common/gvalidator"
)

type tagRepositoryDeleteByID interface {
	DeleteTagByID(ctx context.Context, id string) error
}

type deleteTagByID struct {
	tagRepo tagRepositoryDeleteByID
}

func DeleteByID(tagRepo tagRepositoryDeleteByID) gin.HandlerFunc {
	return func(c *gin.Context) {
		uc := &deleteTagByID{
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

func (uc *deleteTagByID) usecase(ctx context.Context, req *dto.DeleteTagByIDRequest) (*dto.Response, error) {
	err := uc.tagRepo.DeleteTagByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &dto.Response{
		Success: true,
	}, nil
}

func (uc *deleteTagByID) bind(c *gin.Context) (*dto.DeleteTagByIDRequest, error) {
	var req dto.DeleteTagByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (uc *deleteTagByID) validate(req *dto.DeleteTagByIDRequest) error {
	if err := gvalidator.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}

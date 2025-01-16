package tag

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/common/gvalidator"
)

type tagRepositoryDeleteByName interface {
	DeleteTagByName(ctx context.Context, name string) error
}

type deleteTagByName struct {
	tagRepo tagRepositoryDeleteByName
}

func DeleteByName(tagRepo tagRepositoryDeleteByName) gin.HandlerFunc {
	return func(c *gin.Context) {
		uc := &deleteTagByName{
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

func (uc *deleteTagByName) usecase(ctx context.Context, req *dto.DeleteTagByNameRequest) (*dto.Response, error) {
	err := uc.tagRepo.DeleteTagByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return &dto.Response{
		Success: true,
	}, nil
}

func (uc *deleteTagByName) bind(c *gin.Context) (*dto.DeleteTagByNameRequest, error) {
	var req dto.DeleteTagByNameRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (uc *deleteTagByName) validate(req *dto.DeleteTagByNameRequest) error {
	if err := gvalidator.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}

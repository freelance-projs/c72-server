package tag

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/common/gvalidator"
)

type tagRepositoryUpdateTagnameByName interface {
	UpdateTagNameByName(ctx context.Context, oldName, newName string) error
}

type updateTagNameByName struct {
	tagRepo tagRepositoryUpdateTagnameByName
}

func UpdateTagNameByName(tagRepo tagRepositoryUpdateTagnameByName) gin.HandlerFunc {
	return func(c *gin.Context) {
		uc := &updateTagNameByName{
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

func (uc *updateTagNameByName) usecase(ctx context.Context, req *dto.UpdateTagNameRequestByName) (*dto.Response, error) {
	if err := uc.tagRepo.UpdateTagNameByName(ctx, req.OldName, req.NewName); err != nil {
		return nil, err
	}

	return dto.StatusOK(nil), nil
}

func (uc *updateTagNameByName) bind(c *gin.Context) (*dto.UpdateTagNameRequestByName, error) {
	req := dto.UpdateTagNameRequestByName{}
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (uc *updateTagNameByName) validate(req *dto.UpdateTagNameRequestByName) error {
	if err := gvalidator.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}

package tag

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/gvalidator"
	"github.com/xuri/excelize/v2"
)

type tagRepositoryForUpload interface {
	CreateTagInBatches(ctx context.Context, mTags []model.Tag) error
}

type tagMappingUpload struct {
	tagRepo tagRepositoryForUpload
}

func TagMappingUpload(tagRepo tagRepositoryForUpload) gin.HandlerFunc {
	return func(c *gin.Context) {
		uc := &tagMappingUpload{
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

func (t *tagMappingUpload) usecase(ctx context.Context, req *dto.TagMappingUploadRequest) (*dto.Response, error) {
	f, err := excelize.OpenReader(req.Reader)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		_ = f.Close()
	}()
	idToName := t.mapIDToName(f)
	if len(idToName) == 0 {
		return nil, errors.New("no data found in the file")
	}

	mTags := make([]model.Tag, 0, len(idToName))
	for id, name := range idToName {
		mTags = append(mTags, model.Tag{
			ID:   id,
			Name: sql.NullString{String: name, Valid: name != ""},
		})
	}
	if err := t.tagRepo.CreateTagInBatches(ctx, mTags); err != nil {
		return nil, err
	}

	return dto.StatusCreated(nil, "tags"), nil
}

func (t *tagMappingUpload) mapIDToName(f *excelize.File) map[string]string {
	var result = make(map[string]string)

	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		// Get all rows in the sheet
		rows, err := f.GetRows(sheet)
		if err != nil {
			return result
		}

		// Iterate through the rows and print columns A and B
		for _, row := range rows {
			var id, name string
			if len(row) > 0 {
				id = row[0] // Column A
				if len(row) > 1 {
					name = row[1] // Column B
				}
				result[id] = name
			}
		}
	}

	return result
}

func (t *tagMappingUpload) bind(c *gin.Context) (*dto.TagMappingUploadRequest, error) {
	r := c.Request
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return nil, err
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}

	return &dto.TagMappingUploadRequest{
		Reader: file,
	}, nil
}

func (t *tagMappingUpload) validate(req *dto.TagMappingUploadRequest) error {
	if err := gvalidator.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}

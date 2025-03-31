package tagname

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
	"github.com/xuri/excelize/v2"
)

type createByUploadRepo interface {
	CreateTagNameInBatches(ctx context.Context, mTagNames []model.TagName) error
}

type createByUpload struct {
	repo createByUploadRepo
}

func CreateByUpload(repo createByUploadRepo) *createByUpload {
	return &createByUpload{
		repo: repo,
	}
}

func (uc *createByUpload) Usecase(ctx context.Context, req *dto.CreateBatchTagNameRequest) (*ghttp.ResponseBody, error) {
	f, err := excelize.OpenReader(req.Reader)
	if err != nil {
		slog.Error("error opening file", "err", err)
	}
	defer f.Close()

	names, err := uc.getNames(f)
	if err != nil {
		return nil, apperror.ErrBadRequest(err.Error())
	}

	mTagNames := lodash.Map(names, func(name string, _ int) model.TagName {
		return model.TagName{Name: name}
	})

	if err := uc.repo.CreateTagNameInBatches(ctx, mTagNames); err != nil {
		return nil, err
	}

	tagNameDtos := lodash.Map(mTagNames, func(m model.TagName, _ int) dto.TagName {
		return mapper.ToTagNameDto(&m)
	})

	return ghttp.ResponseBodyCreated(tagNameDtos, "tag names are created"), nil
}

func (t *createByUpload) getNames(f *excelize.File) ([]string, error) {
	var names []string

	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		// Get all rows in the sheet
		rows, err := f.GetRows(sheet)
		if err != nil {
			return nil, err
		}

		// Iterate through the rows and print columns A and B
		for _, row := range rows {
			if len(row) > 0 {
				names = append(names, row[0])
			}
		}
	}

	return names, nil
}

func (t *createByUpload) Bind(c *gin.Context) (*dto.CreateBatchTagNameRequest, error) {
	r := c.Request
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return nil, err
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}

	return &dto.CreateBatchTagNameRequest{
		Reader: file,
	}, nil
}

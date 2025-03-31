package department

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
	CreateDepartmentInBatches(ctx context.Context, mDepartments []model.Department) error
}

type createByUpload struct {
	repo createByUploadRepo
}

func CreateByUpload(repo createByUploadRepo) *createByUpload {
	return &createByUpload{
		repo: repo,
	}
}

func (uc *createByUpload) Usecase(ctx context.Context, req *dto.CreateBatchDepartmentRequest) (*ghttp.ResponseBody, error) {
	f, err := excelize.OpenReader(req.Reader)
	if err != nil {
		slog.Error("error opening file", "err", err)
	}
	defer f.Close()

	names, err := uc.getNames(f)
	if err != nil {
		return nil, apperror.ErrBadRequest(err.Error())
	}

	mDepartments := lodash.Map(names, func(name string, _ int) model.Department {
		return model.Department{Name: name}
	})

	if err := uc.repo.CreateDepartmentInBatches(ctx, mDepartments); err != nil {
		return nil, err
	}

	departmentDtos := lodash.Map(mDepartments, func(m model.Department, _ int) dto.Department {
		return mapper.ToDepartmentDto(&m)
	})

	return ghttp.ResponseBodyCreated(departmentDtos, "departments are created"), nil
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

func (t *createByUpload) Bind(c *gin.Context) (*dto.CreateBatchDepartmentRequest, error) {
	r := c.Request
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return nil, err
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}

	return &dto.CreateBatchDepartmentRequest{
		Reader: file,
	}, nil
}

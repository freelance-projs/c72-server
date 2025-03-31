package tagname

import (
	"context"
	"fmt"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type listTagNameRepo interface {
	ListTagNames(ctx context.Context) ([]model.TagName, error)
}

type list struct {
	repo listTagNameRepo
}

func List(repo listTagNameRepo) *list {
	return &list{repo: repo}
}

func (uc *list) Usecase(ctx context.Context, req *dto.ListTagNameRequest) (*ghttp.ResponseBody, error) {
	mTagNames, err := uc.repo.ListTagNames(ctx)

	if err != nil {
		return nil, fmt.Errorf("list tag name error: %w", apperror.GormTranslator(err))
	}

	departments := lodash.Map(mTagNames, func(m model.TagName, _ int) dto.TagName {
		return mapper.ToTagNameDto(&m)
	})

	return ghttp.ResponseBodyOK(departments), nil
}

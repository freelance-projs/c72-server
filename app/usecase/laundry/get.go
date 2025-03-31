package laundry

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

type getLaundryRepo interface {
	GetLaundryByID(ctx context.Context, id int) (*model.Laundry, error)
}

type getLaundry struct {
	repo getLaundryRepo
}

func Get(repo getLaundryRepo) *getLaundry {
	return &getLaundry{repo: repo}
}

func (uc *getLaundry) Usecase(ctx context.Context, req *dto.GetLaundryRequest) (*ghttp.ResponseBody, error) {
	mLaundry, err := uc.repo.GetLaundryByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("get laundry by id error: %w", apperror.GormTranslator(err))
	}

	dtoLaundry := mapper.ToLaundryDto(mLaundry)
	dtoLaundry.Tags = lodash.Map(mLaundry.Tags,
		func(tag model.LaundryTag, i int) dto.LaundryTag {
			return mapper.ToLaundryTagDto(&tag)
		})

	count := lodash.Reduce(mLaundry.Tags, func(acc [2]int, tag model.LaundryTag, _ int) [2]int {
		if tag.Status == model.LaundryStatusWashing {
			acc[0]++
		}
		if tag.Status == model.LaundryStatusReturned {
			acc[1]++
		}
		return acc

	}, [2]int{})
	dtoLaundry.NumWashing = count[0]
	dtoLaundry.NumReturned = count[1]

	return ghttp.ResponseBodyOK(dtoLaundry), nil
}

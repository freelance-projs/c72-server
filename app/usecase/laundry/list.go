package laundry

import (
	"context"
	"fmt"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
	"github.com/ngoctd314/common/qb"
)

type listRepo interface {
	ListLaundries(ctx context.Context, filter qb.Builder) ([]model.Laundry, error)
}

type list struct {
	repo listRepo
}

func List(repo listRepo) *list {
	return &list{repo: repo}
}

func (uc *list) Usecase(ctx context.Context, req *dto.ListLaundryRequest) (*ghttp.ResponseBody, error) {
	builder := qb.New().
		Where(uc.filter(req)).
		Associate(qb.Preload(model.Laundry{}.TagsRelation()))
	mLaundries, err := uc.repo.ListLaundries(ctx, builder)
	if err != nil {
		return nil, fmt.Errorf("list laundries error: %w", apperror.GormTranslator(err))
	}

	dtoLaundry := lodash.Map(mLaundries, func(m model.Laundry, _ int) dto.Laundry {
		laundryDTO := mapper.ToLaundryDto(&m)

		numReturned := 0
		for _, tag := range m.Tags {
			if tag.Status == model.LaundryStatusReturned {
				numReturned++
			}
		}

		laundryDTO.NumReturned = numReturned

		return laundryDTO
	})

	return ghttp.ResponseBodyOK(dtoLaundry), nil
}

const weakDuration = time.Hour * 24 * 7

func (uc *list) filter(req *dto.ListLaundryRequest) *qb.Cond {
	mLaundryColumns := model.Laundry{}.Columns()
	filters := []*qb.Cond{
		qb.Gte(mLaundryColumns.CreatedAt, time.Unix(*req.From, 0)),
		qb.Lte(mLaundryColumns.CreatedAt, time.Unix(*req.To, 0)),
	}
	if req.Completed != nil {
		if *req.Completed {
			filters = append(filters, qb.Eq(mLaundryColumns.NumWashing, 0))
		} else {
			filters = append(filters, qb.NotEq(mLaundryColumns.NumWashing, 0))
		}
	}

	return qb.And(filters...)
}

func (uc *list) Validate(ctx context.Context, req *dto.ListLaundryRequest) error {
	now := time.Now()

	if req.From == nil && req.To == nil {
		to := now.Unix()
		req.To = &to
		from := now.Add(-weakDuration).Unix()
		req.From = &from
	}

	if req.To == nil {
		to := now.Unix()
		req.To = &to
	}

	if req.From == nil {
		from := time.Unix(*req.To, 0).Add(-weakDuration).Unix()
		req.From = &from
	}

	return nil
}

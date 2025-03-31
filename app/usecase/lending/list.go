package lending

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
	ListLending(ctx context.Context, filter qb.Builder) ([]model.Lending, error)
}

type list struct {
	repo listRepo
}

func List(repo listRepo) *list {
	return &list{repo: repo}
}

func (uc *list) Usecase(ctx context.Context, req *dto.ListLendingRequest) (*ghttp.ResponseBody, error) {
	builder := qb.New().
		Where(uc.filter(req)).
		Associate(qb.Preload(model.Lending{}.TagsRelation()))

	mLendings, err := uc.repo.ListLending(ctx, builder)
	if err != nil {
		return nil, fmt.Errorf("list lending error: %w", apperror.GormTranslator(err))
	}

	dtoLending := lodash.Map(mLendings, func(m model.Lending, _ int) dto.Lending {
		lendingDTO := mapper.ToLendingDTO(&m)
		numReturned := 0
		for _, tag := range m.Tags {
			if tag.Status == model.LendingStatusReturned {
				numReturned++
			}
		}

		lendingDTO.NumReturned = numReturned

		return lendingDTO
	})

	return ghttp.ResponseBodyOK(dtoLending), nil
}

const weakDuration = time.Hour * 24 * 7

func (uc *list) filter(req *dto.ListLendingRequest) *qb.Cond {
	mLendingColumns := model.Lending{}.Columns()
	filters := []*qb.Cond{
		qb.Gte(mLendingColumns.CreatedAt, time.Unix(*req.From, 0)),
		qb.Lte(mLendingColumns.CreatedAt, time.Unix(*req.To, 0)),
	}
	if req.Completed != nil {
		if *req.Completed {
			filters = append(filters, qb.Eq(mLendingColumns.NumLending, 0))
		} else {
			filters = append(filters, qb.NotEq(mLendingColumns.NumLending, 0))
		}
	}

	return qb.And(filters...)
}

func (uc *list) Validate(ctx context.Context, req *dto.ListLendingRequest) error {
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

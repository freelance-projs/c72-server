package txlog

import (
	"context"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
	"github.com/ngoctd314/common/qb"
)

type listDept struct {
	repo *repository.Laundry
}

func ListDept(repo *repository.Laundry) *listDept {
	return &listDept{repo: repo}
}

func (uc *listDept) Usecase(ctx context.Context, req *dto.ListTxLogDeptRequest) (*ghttp.ResponseBody, error) {
	filter := qb.New().Where(uc.filter(req))
	mTxLogs, err := uc.repo.ListTxLogDept(ctx, filter)
	if err != nil {
		return nil, err
	}

	txLogDtos := lodash.Map(mTxLogs, func(m model.TxLogDepartment, _ int) dto.TxLogDepartment {
		return mapper.ToTxLogDept(&m)
	})

	return ghttp.ResponseBodyOK(txLogDtos), nil
}

func (uc *listDept) filter(req *dto.ListTxLogDeptRequest) *qb.Cond {
	filters := []*qb.Cond{
		qb.Gte("created_at", time.Unix(*req.From, 0)),
		qb.Lte("created_at", time.Unix(*req.To, 0)),
	}

	return qb.And(filters...)
}

const weakDuration = time.Hour * 24 * 7

func (uc *listDept) Validate(ctx context.Context, req *dto.ListTxLogDeptRequest) error {
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

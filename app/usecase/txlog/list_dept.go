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
)

type listDept struct {
	repo *repository.Repository
}

func ListDept(repo *repository.Repository) *listDept {
	return &listDept{repo: repo}
}

func (uc *listDept) Usecase(ctx context.Context, req *dto.ListTxLogRequest) (*ghttp.ResponseBody, error) {
	from := time.Unix(*req.From, 0)
	to := time.Unix(*req.To, 0)
	mLendingStats, err := uc.repo.ListTxLogDepartments(ctx, from, to)
	if err != nil {
		return nil, err
	}

	dtoLendingStats := lodash.Map(mLendingStats, func(stat model.LendingStat, _ int) dto.TxLogDepartment {
		return mapper.ToTxLogDepartmentDto(&stat)
	})

	return ghttp.ResponseBodyOK(dtoLendingStats), nil
}

const weakDuration = time.Hour * 24 * 7

func (uc *listDept) Validate(ctx context.Context, req *dto.ListTxLogRequest) error {
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

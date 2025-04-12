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

type listCompany struct {
	repo *repository.Repository
}

func ListCompany(repo *repository.Repository) *listCompany {
	return &listCompany{repo: repo}
}

func (uc *listCompany) Usecase(ctx context.Context, req *dto.ListTxLogDeptRequest) (*ghttp.ResponseBody, error) {
	from := time.Unix(*req.From, 0)
	to := time.Unix(*req.To, 0)
	mTxLogs, err := uc.repo.ListTxLogCompany(ctx, from, to)
	if err != nil {
		return nil, err
	}

	dtoLendingStats := lodash.Map(mTxLogs, func(stat model.TxLogCompany, _ int) dto.TxLogCompany {
		return mapper.ToLogCompanyDto(&stat)
	})

	return ghttp.ResponseBodyOK(dtoLendingStats), nil
}

func (uc *listCompany) Validate(ctx context.Context, req *dto.ListTxLogDeptRequest) error {
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

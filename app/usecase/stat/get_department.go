package stat

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type getDepartment struct {
	repo *repository.Repository
}

func GetDepartment(repo *repository.Repository) *getDepartment {
	return &getDepartment{repo: repo}
}

func (uc *getDepartment) Usecase(ctx context.Context, req *dto.GetDepartmentStatRequest) (*ghttp.ResponseBody, error) {
	mStats, err := uc.repo.GetDepartmentStat(ctx, req.Department)
	if err != nil {
		return nil, err
	}

	numExported, numReturned := 0, 0
	trackings := lodash.Map(mStats, func(stat model.LendingStat, _ int) dto.TagTracking {
		numExported += stat.Lending
		numReturned += stat.Returned
		return dto.TagTracking{
			Name:     stat.TagName,
			Exported: stat.Lending,
			Returned: stat.Returned,
		}
	})

	resp := dto.DepartmentDetailStat{
		Department: req.Department,
		Exported:   numExported,
		Returned:   numReturned,
		Trackings:  trackings,
	}

	return ghttp.ResponseBodyOK(resp), nil
}

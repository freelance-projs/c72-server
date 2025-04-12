package stat

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type getTag struct {
	repo *repository.Repository
}

func GetTag(repo *repository.Repository) *getTag {
	return &getTag{repo: repo}
}

func (uc *getTag) Usecase(ctx context.Context, req *dto.GetTagStatRequest) (*ghttp.ResponseBody, error) {
	mStatsDept, err := uc.repo.GetLendingTagStat(ctx, req.TagName)
	if err != nil {
		return nil, err
	}

	mStatsCompany, err := uc.repo.GetWashingTagStat(ctx, req.TagName)
	if err != nil {
		return nil, err
	}

	lending, lendingReturned := 0, 0
	departmentTracking := lodash.Map(mStatsDept, func(stat model.LendingStat, _ int) dto.DepartmentTracking {
		lending += stat.Lending
		lendingReturned += stat.Returned
		return dto.DepartmentTracking{
			Name:     stat.Department,
			Exported: stat.Lending,
			Returned: stat.Returned,
		}
	})
	washing, washingReturned := 0, 0
	companyTracking := lodash.Map(mStatsCompany, func(stat model.WashingStat, _ int) dto.CompanyTracking {
		washing += stat.Washing
		washingReturned += stat.Returned
		return dto.CompanyTracking{
			Name:     stat.Company,
			Exported: stat.Washing,
			Returned: stat.Returned,
		}
	})

	resp := dto.TagStatDetail{
		TagName:         req.TagName,
		Lending:         lending,
		LendingReturned: lendingReturned,
		Washing:         washing,
		WashingReturned: washingReturned,
		Departments:     departmentTracking,
		Companies:       companyTracking,
	}

	return ghttp.ResponseBodyOK(resp), nil
}

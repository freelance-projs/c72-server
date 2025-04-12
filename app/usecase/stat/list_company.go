package stat

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/net/ghttp"
)

type listCompany struct {
	repo *repository.Repository
}

func ListCompany(repo *repository.Repository) *listCompany {
	return &listCompany{repo: repo}
}

func (uc *listCompany) Usecase(ctx context.Context, req *dto.ListCompaniesRequest) (*ghttp.ResponseBody, error) {
	mStats, err := uc.repo.ListCompanyStats(ctx)
	if err != nil {
		return nil, err
	}
	dtoStats := make([]dto.CompanyStat, 0, len(mStats))
	m := make(map[string]int)
	for _, stat := range mStats {
		if _, ok := m[stat.Company]; !ok {
			dtoStats = append(dtoStats, dto.CompanyStat{
				Company: stat.Company,
			})
			m[stat.Company] = len(dtoStats) - 1
		}

		dtoStats[m[stat.Company]].Exported += stat.Washing
		dtoStats[m[stat.Company]].Returned += stat.Returned
	}

	return ghttp.ResponseBodyOK(dtoStats), nil
}

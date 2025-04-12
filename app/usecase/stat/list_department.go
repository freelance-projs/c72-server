package stat

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/net/ghttp"
)

type listDepartment struct {
	repo *repository.Repository
}

func ListDepartment(repo *repository.Repository) *listDepartment {
	return &listDepartment{repo: repo}
}

func (uc *listDepartment) Usecase(ctx context.Context, req *dto.ListDepartmentStatRequest) (*ghttp.ResponseBody, error) {
	mStats, err := uc.repo.ListDepartmentStats(ctx)
	if err != nil {
		return nil, err
	}

	dtoStats := make([]dto.DepartmentStat, 0, len(mStats))
	m := make(map[string]int)
	for _, stat := range mStats {
		if _, ok := m[stat.Department]; !ok {
			dtoStats = append(dtoStats, dto.DepartmentStat{
				Department: stat.Department,
			})
			m[stat.Department] = len(dtoStats) - 1
		}

		dtoStats[m[stat.Department]].Exported += stat.Lending
		dtoStats[m[stat.Department]].Returned += stat.Returned
	}

	return ghttp.ResponseBodyOK(dtoStats), nil
}

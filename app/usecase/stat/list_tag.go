package stat

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/net/ghttp"
)

type listTag struct {
	repo *repository.Repository
}

func ListTag(repo *repository.Repository) *listTag {
	return &listTag{repo: repo}
}

func (uc *listTag) Usecase(ctx context.Context, req *dto.ListTagStatRequest) (*ghttp.ResponseBody, error) {
	lendingStats, err := uc.repo.ListLendingTagStats(ctx)
	if err != nil {
		return nil, err
	}
	washingStats, err := uc.repo.ListWashingTagStats(ctx)
	if err != nil {
		return nil, err
	}

	dept := make(map[string][2]int)
	for _, stat := range lendingStats {
		dept[stat.TagName] = [2]int{stat.Lending, stat.Returned}
	}

	company := make(map[string][2]int)
	for _, stat := range washingStats {
		company[stat.TagName] = [2]int{stat.Washing, stat.Returned}
	}

	m := make(map[string]int)
	dtoStats := make([]dto.TagStat, 0)
	for _, stat := range lendingStats {
		m[stat.TagName] = len(dtoStats)
		dtoStats = append(dtoStats, dto.TagStat{
			TagName:         stat.TagName,
			Lending:         stat.Lending,
			LendingReturned: stat.Returned,
		})
	}
	for _, stat := range washingStats {
		if _, ok := m[stat.TagName]; !ok {
			m[stat.TagName] = len(dtoStats)
			dtoStats = append(dtoStats, dto.TagStat{
				TagName:         stat.TagName,
				Washing:         stat.Washing,
				WashingReturned: stat.Returned,
			})
		} else {
			dtoStats[m[stat.TagName]].Washing += stat.Washing
			dtoStats[m[stat.TagName]].WashingReturned += stat.Returned
		}
	}

	return ghttp.ResponseBodyOK(dtoStats), nil
}

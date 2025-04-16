package stat

import (
	"context"
	"log/slog"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/net/ghttp"
)

type listTag struct {
	repo     *repository.Repository
	sheetSvc *sheetService
	sheetID  string
}

func ListTag(repo *repository.Repository) *listTag {
	uc := &listTag{
		repo:     repo,
		sheetSvc: newSheetService(),
	}

	setting, err := repo.GetSetting(context.Background())
	if err == nil {
		uc.sheetID = setting.TxLogSheetID
	} else {
		uc.sheetID = "1xgd39AuKdQKnyOJO63W7Y3KueUVyoBdsYskhRMpOKW4"
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			setting, err := repo.GetSetting(context.Background())
			slog.Info("updating stat sheet id")
			if err == nil {
				uc.sheetID = setting.TxLogSheetID
			}
			ticker.Reset(time.Minute)
		}
	}()
	return uc
}

func (uc *listTag) Usecase(ctx context.Context, req *dto.ListTagStatRequest) (*ghttp.ResponseBody, error) {
	from := time.Unix(*req.From, 0)
	to := time.Unix(*req.To, 0)

	lendingStats, err := uc.repo.ListLendingTagStats(ctx, from, to)
	if err != nil {
		return nil, err
	}
	washingStats, err := uc.repo.ListWashingTagStats(ctx, from, to)
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

	go func() {
		sheetCols := make([]any, 0, len(dtoStats))
		for _, stat := range dtoStats {
			sheetCols = append(sheetCols, stat)
		}

		now := time.Now()
		sheetName := "Thẻ " + time.Now().Format("2006-01-02")
		if err := uc.sheetSvc.insert(uc.sheetID, sheetName, sheetCols); err != nil {
			slog.Error("error inserting data to sheet", "err", err)
		}
		slog.Info("insert data to sheet successfully", "sheetID", uc.sheetID, "since", time.Since(now).Seconds())
	}()

	return ghttp.ResponseBodyOK(dtoStats), nil
}

func (uc *listTag) Validate(ctx context.Context, req *dto.ListTagStatRequest) error {
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

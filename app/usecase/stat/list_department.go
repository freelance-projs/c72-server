package stat

import (
	"context"
	"log/slog"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/net/ghttp"
)

type listDepartment struct {
	repo     *repository.Repository
	sheetSvc *sheetService
	sheetID  string
}

func ListDepartment(repo *repository.Repository) *listDepartment {
	uc := &listDepartment{
		repo:     repo,
		sheetSvc: newSheetService(),
	}

	setting, err := repo.GetSetting(context.Background())
	if err == nil {
		uc.sheetID = setting.ReportSheetID
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

func (uc *listDepartment) Usecase(ctx context.Context, req *dto.ListDepartmentStatRequest) (*ghttp.ResponseBody, error) {
	from := time.Unix(*req.From, 0)
	to := time.Unix(*req.To, 0)
	mStats, err := uc.repo.ListDepartmentStats(ctx, from, to)
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

	if len(dtoStats) > 0 {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("error inserting data to sheet", "err", r)
				}
			}()
			sheetCols := make([]any, 0, len(dtoStats))
			for _, stat := range dtoStats {
				sheetCols = append(sheetCols, stat)
			}

			now := time.Now()
			sheetName := "Nội bộ " + time.Now().Format("2006-01-02")
			if err := uc.sheetSvc.insert(uc.sheetID, sheetName, sheetCols); err != nil {
				slog.Error("error inserting data to sheet", "err", err)
			}
			slog.Info("insert data to sheet successfully", "sheetID", uc.sheetID, "since", time.Since(now).Seconds())
		}()
	}

	return ghttp.ResponseBodyOK(dtoStats), nil
}

const weakDuration = time.Hour * 24 * 7

func (uc *listDepartment) Validate(ctx context.Context, req *dto.ListDepartmentStatRequest) error {
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

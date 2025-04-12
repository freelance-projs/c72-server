package stat

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/net/ghttp"
)

type getCompany struct {
	repo *repository.Repository
}

func GetCompany(repo *repository.Repository) *getCompany {
	return &getCompany{repo: repo}
}

func (uc *getCompany) Usecase(ctx context.Context, req *dto.GetCompanyStatRequest) (*ghttp.ResponseBody, error) {
	mStats, err := uc.repo.GetCompanyStat(ctx, req.Company)
	if err != nil {
		return nil, err
	}

	numExported, numReturned := 0, 0
	trackings := make([]dto.TagTracking, 0, len(mStats))
	for _, stat := range mStats {
		numExported += stat.Washing
		numReturned += stat.Returned
		trackings = append(trackings, dto.TagTracking{
			Name:     stat.TagName,
			Exported: stat.Washing,
			Returned: stat.Returned,
		})
	}

	resp := dto.CompanyDetailStat{
		Company:   req.Company,
		Exported:  numExported,
		Returned:  numReturned,
		Trackings: trackings,
	}

	return ghttp.ResponseBodyOK(resp), nil
}

package txlog

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

func (uc *getCompany) Usecase(ctx context.Context, req *dto.GetTxLogCompanyRequest) (*ghttp.ResponseBody, error) {
	mStats, err := uc.repo.GetTxLogCompanyByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	sumWashing, sumReturned := 0, 0
	statDetails := make([]dto.TxLogDetail, 0, len(mStats))

	type aggKey struct {
		company   string
		createdAt int
	}
	aggTagTracking := make(map[aggKey][]dto.TagTracking)
	for _, v := range mStats {
		sumWashing += int(v.Washing)
		sumReturned += int(v.Returned)
		key := aggKey{
			company:   v.Company,
			createdAt: int(v.CreatedAt.Unix()),
		}
		if _, ok := aggTagTracking[key]; !ok {
			statDetails = append(statDetails, dto.TxLogDetail{
				Entity:    v.Company,
				Action:    v.Action.String(),
				Actor:     v.Actor,
				CreatedAt: v.CreatedAt,
			})
		}
		aggTagTracking[key] = append(aggTagTracking[key], dto.TagTracking{
			Name:     v.TagName,
			Exported: v.Washing,
			Returned: v.Returned,
		})
	}

	for i := range statDetails {
		statDetails[i].Tracking = aggTagTracking[aggKey{
			company:   statDetails[i].Entity,
			createdAt: int(statDetails[i].CreatedAt.Unix()),
		}]
	}

	dtoDeptstat := dto.TxLog{
		ID:       req.ID,
		Exported: sumWashing,
		Returned: sumReturned,
		Details:  statDetails,
	}

	return ghttp.ResponseBodyOK(dtoDeptstat), nil

}

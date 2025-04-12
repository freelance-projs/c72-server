package txlog

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/net/ghttp"
)

type getDept struct {
	repo *repository.Repository
}

func GetDept(repo *repository.Repository) *getDept {
	return &getDept{repo: repo}
}

func (uc *getDept) Usecase(ctx context.Context, req *dto.GetTxLogDepartmentRequest) (*ghttp.ResponseBody, error) {
	mStats, err := uc.repo.GetTxLogDepartmentByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	sumLending, sumReturned := 0, 0
	statDetails := make([]dto.TxLogDetail, 0, len(mStats))

	type aggKey struct {
		department string
		createdAt  int
	}
	aggTagTracking := make(map[aggKey][]dto.TagTracking)
	for _, v := range mStats {
		sumLending += int(v.Lending)
		sumReturned += int(v.Returned)
		key := aggKey{
			department: v.Department,
			createdAt:  int(v.CreatedAt.Unix()),
		}
		if _, ok := aggTagTracking[key]; !ok {
			statDetails = append(statDetails, dto.TxLogDetail{
				Entity:    v.Department,
				Action:    v.Action.String(),
				Actor:     v.Actor,
				CreatedAt: v.CreatedAt,
			})
		}
		aggTagTracking[key] = append(aggTagTracking[key], dto.TagTracking{
			Name:     v.TagName,
			Exported: v.Lending,
			Returned: v.Returned,
		})
	}

	for i := range statDetails {
		statDetails[i].Tracking = aggTagTracking[aggKey{
			department: statDetails[i].Entity,
			createdAt:  int(statDetails[i].CreatedAt.Unix()),
		}]
	}

	dtoDeptstat := dto.TxLog{
		ID:       req.ID,
		Exported: sumLending,
		Returned: sumReturned,
		Details:  statDetails,
	}

	return ghttp.ResponseBodyOK(dtoDeptstat), nil
}

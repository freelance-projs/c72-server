package txlog

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type getCompany struct {
	repo *repository.Laundry
}

func GetCompany(repo *repository.Laundry) *getCompany {
	return &getCompany{repo: repo}
}

func (uc *getCompany) Usecase(ctx context.Context, req *dto.GetTxLogDeptRequest) (*ghttp.ResponseBody, error) {
	mTxLog, err := uc.repo.GetTxLogCompany(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	details := make([]dto.TxLogDetail, 0, len(mTxLog.Details))
	for _, v := range mTxLog.Details {
		details = append(details, dto.TxLogDetail{
			Action: v.Action,
			Tracking: lodash.Map(v.Tracking, func(tracking model.TxLogTracking, _ int) dto.TxLogTracking {
				return dto.TxLogTracking{
					Name:  tracking.Name,
					Count: tracking.Count,
				}
			}),
			CreatedAt: v.CreatedAt,
		})
	}

	return ghttp.ResponseBodyOK(dto.TxLogCompany{
		ID:          mTxLog.ID,
		Company:     mTxLog.Overview.Actor,
		NumWashing:  int(mTxLog.Overview.TotalTags - mTxLog.Overview.Returned),
		NumReturned: int(mTxLog.Overview.Returned),
		Details:     details,
		CreatedAt:   mTxLog.CreatedAt,
	}), nil
}

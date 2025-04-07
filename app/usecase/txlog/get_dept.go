package txlog

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type getDept struct {
	repo *repository.Laundry
}

func GetDept(repo *repository.Laundry) *getDept {
	return &getDept{repo: repo}
}

func (uc *getDept) Usecase(ctx context.Context, req *dto.GetTxLogDeptRequest) (*ghttp.ResponseBody, error) {
	mTxLog, err := uc.repo.GetTxLogDept(ctx, req.ID)
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

	return ghttp.ResponseBodyOK(dto.TxLogDepartment{
		ID:          mTxLog.ID,
		Department:  mTxLog.Overview.Actor,
		NumLending:  int(mTxLog.Overview.TotalTags - mTxLog.Overview.Returned),
		NumReturned: int(mTxLog.Overview.Returned),
		Details:     details,
		CreatedAt:   mTxLog.CreatedAt,
	}), nil
}

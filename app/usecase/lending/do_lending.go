package lending

import (
	"context"
	"fmt"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/service"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type doLendingRepo interface {
	CreateLending(ctx context.Context, mLending *model.Lending, tagIDs []string) error
	GetTagsByIDs(ctx context.Context, ids []string) ([]model.Tag, error)
}

type doLending struct {
	repo doLendingRepo
}

func DoLending(repo doLendingRepo) *doLending {
	return &doLending{
		repo: repo,
	}
}

func (uc *doLending) Usecase(ctx context.Context, req *dto.DoLendingRequest) (*ghttp.ResponseBody, error) {
	mTags, err := uc.repo.GetTagsByIDs(ctx, req.TagIDs)
	if err != nil {
		return nil, fmt.Errorf("get tags by ids error: %w", apperror.GormTranslator(err))
	}
	now := time.Now()
	mTags = lodash.Filter(mTags, func(m model.Tag, _ int) bool {
		if !m.LastUsed.Valid {
			return true
		}
		return m.LastUsed.Time.Add(service.SystemSetting.GetTxLockTime()).Before(now)
	})

	if len(mTags) == 0 {
		return nil, apperror.New("các tag đang bị khoá, vui lòng thử lại sau")
	}

	tagIDs := lodash.Map(mTags, func(m model.Tag, _ int) string {
		return m.ID
	})

	mLending := &model.Lending{
		Department: req.Department,
		NumLending: len(tagIDs),
	}
	if err := uc.repo.CreateLending(ctx, mLending, tagIDs); err != nil {
		return nil, fmt.Errorf("create lending error: %w", apperror.GormTranslator(err))
	}

	lending := mapper.ToLendingDTO(mLending)

	return ghttp.ResponseBodyCreated(lending, "Tạo giao dịch mượn đồ thành công"), nil
}

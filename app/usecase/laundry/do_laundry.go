package laundry

// import (
// 	"context"
// 	"fmt"
// 	"time"
//
// 	"github.com/ngoctd314/c72-api-server/pkg/dto"
// 	"github.com/ngoctd314/c72-api-server/pkg/mapper"
// 	"github.com/ngoctd314/c72-api-server/pkg/model"
// 	"github.com/ngoctd314/c72-api-server/pkg/service"
// 	"github.com/ngoctd314/common/apperror"
// 	"github.com/ngoctd314/common/lodash"
// 	"github.com/ngoctd314/common/net/ghttp"
// )
//
// type doLaundryRepo interface {
// 	CreateLaundry(ctx context.Context, mLaundry *model.Laundry, tagIDs []string) error
// 	GetTagsByIDs(ctx context.Context, ids []string) ([]model.Tag, error)
// }
//
// type doLaundry struct {
// 	repo doLaundryRepo
// }
//
// func DoLaundry(repo doLaundryRepo) *doLaundry {
// 	return &doLaundry{
// 		repo: repo,
// 	}
// }
//
// func (uc *doLaundry) Usecase(ctx context.Context, req *dto.DoLaundryRequest) (*ghttp.ResponseBody, error) {
// 	mTags, err := uc.repo.GetTagsByIDs(ctx, req.TagIDs)
// 	if err != nil {
// 		return nil, fmt.Errorf("get tags by ids error: %w", apperror.GormTranslator(err))
// 	}
// 	now := time.Now()
// 	mTags = lodash.Filter(mTags, func(m model.Tag, _ int) bool {
// 		if !m.LastWashing.Valid {
// 			return true
// 		}
// 		return m.LastWashing.Time.Add(service.SystemSetting.GetTxLockTime()).Before(now)
// 	})
//
// 	if len(mTags) == 0 {
// 		return nil, apperror.New("các tag đang bị khoá, vui lòng thử lại sau")
// 	}
//
// 	tagIDs := lodash.Map(mTags, func(m model.Tag, _ int) string {
// 		return m.ID
// 	})
// 	mLaundry := &model.Laundry{
// 		Name:       req.Company,
// 		NumWashing: len(tagIDs),
// 	}
//
// 	if err := uc.repo.CreateLaundry(ctx, mLaundry, tagIDs); err != nil {
// 		return nil, fmt.Errorf("create laundry error: %w", apperror.GormTranslator(err))
// 	}
//
// 	laundryDto := mapper.ToLaundryDto(mLaundry)
//
// 	return ghttp.ResponseBodyCreated(laundryDto, "Tạo giao dịch giặt đồ thành công"), nil
// }

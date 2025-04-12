package tag

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type tagRepoListByID interface {
	GetTagsByIDs(ctx context.Context, ids []string) ([]model.Tag, error)
	GetActiveTag(ctx context.Context, tagIDs []string, action string) ([]model.Tag, error)
}

type listTagByID struct {
	repo *repository.Repository
}

func GetActiveTagsByIDs(repo *repository.Repository) *listTagByID {
	return &listTagByID{
		repo: repo,
	}
}

func (uc *listTagByID) Usecase(ctx context.Context, req *dto.GetActiveTagsRequest) (*ghttp.ResponseBody, error) {
	mTags, err := uc.repo.GetActiveTags(ctx, req.Action, req.IDs)
	if err != nil {
		return nil, err
	}

	// now := time.Now()
	// mTags = lodash.Filter(mTags, func(m model.Tag, _ int) bool {
	// 	switch req.Type {
	// 	case "last_used":
	// 		if !m.LastUsed.Valid {
	// 			return true
	// 		}
	// 		return m.LastUsed.Time.Add(service.SystemSetting.GetTxLockTime()).Before(now)
	// 	case "last_washing":
	// 		if !m.LastWashing.Valid {
	// 			return true
	// 		}
	// 		return m.LastWashing.Time.Add(service.SystemSetting.GetTxLockTime()).Before(now)
	// 	default:
	// 		return true
	// 	}
	// })

	tagDtos := lodash.Map(mTags, func(mTag model.Tag, _ int) dto.Tag {
		return mapper.ToTagDto(&mTag)
	})

	return ghttp.ResponseBodyOK(tagDtos), nil
}

package tag

import (
	"context"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/service"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type tagRepoListByID interface {
	GetTagsByIDs(ctx context.Context, ids []string) ([]model.Tag, error)
}

type listTagByID struct {
	tagRepo tagRepoListByID
}

func ListTagByID(tagRepo tagRepoListByID) *listTagByID {
	return &listTagByID{
		tagRepo: tagRepo,
	}
}

func (uc *listTagByID) Usecase(ctx context.Context, req *dto.ListTagByIDRequest) (*ghttp.ResponseBody, error) {
	mTags, err := uc.tagRepo.GetTagsByIDs(ctx, req.IDs)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	mTags = lodash.Filter(mTags, func(m model.Tag, _ int) bool {
		switch req.Type {
		case "last_used":
			if !m.LastUsed.Valid {
				return true
			}
			return m.LastUsed.Time.Add(service.SystemSetting.GetTxLockTime()).Before(now)
		case "last_washing":
			if !m.LastWashing.Valid {
				return true
			}
			return m.LastWashing.Time.Add(service.SystemSetting.GetTxLockTime()).Before(now)
		default:
			return true
		}
	})

	tagDtos := lodash.Map(mTags, func(mTag model.Tag, _ int) dto.Tag {
		return mapper.ToTagDto(&mTag)
	})

	return ghttp.ResponseBodyOK(tagDtos), nil
}

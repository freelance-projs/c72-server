package lending

import (
	"context"
	"fmt"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/mapper"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/lodash"
	"github.com/ngoctd314/common/net/ghttp"
)

type getLendingTagRepo interface {
	GetLendingByID(ctx context.Context, id int) (*model.Lending, error)
}

type getLending struct {
	repo getLendingTagRepo
}

func GetLendingTag(repo getLendingTagRepo) *getLending {
	return &getLending{repo: repo}
}

func (uc *getLending) Usecase(ctx context.Context, req *dto.GetLendingRequest) (*ghttp.ResponseBody, error) {
	mLending, err := uc.repo.GetLendingByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("get lending by id error: %w", apperror.GormTranslator(err))
	}

	dtoLending := mapper.ToLendingDTO(mLending)
	dtoLending.Tags = lodash.Map(mLending.Tags,
		func(tag model.LendingTag, i int) dto.LendingTag {
			return mapper.ToLendingTagDTO(&tag)
		})

	count := lodash.Reduce(mLending.Tags, func(acc [2]int, tag model.LendingTag, _ int) [2]int {
		if tag.Status == model.LendingStatusLending {
			acc[0]++
		}
		if tag.Status == model.LendingStatusReturned {
			acc[1]++
		}
		return acc
	}, [2]int{})
	dtoLending.NumLending = count[0]
	dtoLending.NumReturned = count[1]

	return ghttp.ResponseBodyOK(dtoLending), nil
}

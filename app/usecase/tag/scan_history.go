package tag

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/common/gvalidator"
)

type tagRepositoryScanHistory interface {
	GetTagsScanHistories(ctx context.Context, from, to *time.Time) ([]model.TagScanHistory, error)
	GetTagsByIDs(ctx context.Context, ids []string) ([]model.Tag, error)
}

type tagScanHistory struct {
	tagRepo tagRepositoryScanHistory
}

func NewTagScanHistory(tagRepo tagRepositoryScanHistory) gin.HandlerFunc {
	return func(c *gin.Context) {
		uc := &tagScanHistory{
			tagRepo: tagRepo,
		}

		req, bindErr := uc.bind(c)
		if bindErr != nil {
			dto.JSONFail(c, bindErr)
			return
		}

		if validateErr := uc.validate(req); validateErr != nil {
			dto.JSONFail(c, validateErr)
			return
		}

		resp, err := uc.usecase(c.Request.Context(), req)
		if err != nil {
			dto.JSONFail(c, err)
			return
		}

		dto.JSONSuccess(c, resp)
	}
}

func (uc *tagScanHistory) usecase(ctx context.Context, req *dto.TagScanHistoryRequest) (*dto.Response, error) {
	mScanHistories, err := uc.tagRepo.GetTagsScanHistories(ctx, req.From, req.To)
	if err != nil {
		return nil, err
	}
	if len(mScanHistories) == 0 {
		return dto.StatusOK(nil), nil
	}

	uniqueTagIDs := make(map[string]struct{})
	tagIDs := make([]string, 0, len(mScanHistories))
	for _, v := range mScanHistories {
		if _, ok := uniqueTagIDs[v.TagID]; !ok {
			tagIDs = append(tagIDs, v.TagID)
		}
		uniqueTagIDs[v.TagID] = struct{}{}
	}

	mTags, err := uc.tagRepo.GetTagsByIDs(ctx, tagIDs)
	if err != nil {
		return nil, err
	}
	if len(mTags) == 0 {
		return dto.StatusOK(nil), nil
	}

	resp := uc.buildResponse(mScanHistories, mTags)

	return dto.StatusOK(resp), nil
}

func (uc *tagScanHistory) buildResponse(mScanHistories []model.TagScanHistory, mTags []model.Tag) []dto.TagScanHistoryResponse {
	mTagIDToTag := make(map[string]model.Tag)
	for _, v := range mTags {
		mTagIDToTag[v.ID] = v
	}

	// result
	var resp []dto.TagScanHistoryResponse
	tracking := make(map[[2]string]int)

	currentDay := mScanHistories[0].CreatedAt
	currentHistory := dto.TagScanHistoryResponse{
		Day: currentDay.Format("2006-01-02"),
	}
	// map tag to day and scan time
	for _, v := range mScanHistories {
		if v.CreatedAt.Day() != currentDay.Day() {
			currentDay = v.CreatedAt
			resp = append(resp, currentHistory)
			currentHistory = dto.TagScanHistoryResponse{
				Day: currentDay.Format("2006-01-02"),
			}
		}
		tagName := mTagIDToTag[v.TagID].Name.String
		if tagName == "" {
			tagName = v.TagID
		}
		k := [2]string{v.CreatedAt.String(), tagName}
		if _, ok := tracking[k]; !ok {
			tracking[k] = len(currentHistory.Histories)
			currentHistory.Histories = append(currentHistory.Histories, dto.ScanHistoryUnit{
				TagName:   tagName,
				Count:     1,
				CreatedAt: v.CreatedAt,
			})
		} else {
			currentHistory.Histories[tracking[k]].Count++
		}
	}
	return append(resp, currentHistory)
}

func (uc *tagScanHistory) bind(c *gin.Context) (*dto.TagScanHistoryRequest, error) {
	var req dto.TagScanHistoryRequest

	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	if req.From == nil {
		req.From = &time.Time{}
	}

	if req.To == nil {
		now := time.Now()
		req.To = &now
	}

	if req.To.Before(*req.From) {
		req.From, req.To = req.To, req.From
	}

	return &req, nil
}

func (uc *tagScanHistory) validate(req *dto.TagScanHistoryRequest) error {
	if err := gvalidator.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}

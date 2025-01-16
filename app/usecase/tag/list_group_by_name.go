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
	GetTags(ctx context.Context, from, to *time.Time) ([]model.Tag, error)
}

type listGroupByName struct {
	tagRepo tagRepositoryScanHistory
}

func NewTagScanHistory(tagRepo tagRepositoryScanHistory) gin.HandlerFunc {
	return func(c *gin.Context) {
		uc := &listGroupByName{
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

func (uc *listGroupByName) usecase(ctx context.Context, req *dto.TagScanHistoryRequest) (*dto.Response, error) {
	mTags, err := uc.tagRepo.GetTags(ctx, req.From, req.To)
	if err != nil {
		return nil, err
	}
	if len(mTags) == 0 {
		return dto.StatusOK(nil), nil
	}

	resp := uc.buildResponse(mTags)

	return dto.StatusOK(resp), nil
}

func (uc *listGroupByName) buildResponse(mTags []model.Tag) []dto.TagScanHistoryResponseV2 {
	mTagIDToTag := make(map[string]model.Tag)
	for _, v := range mTags {
		mTagIDToTag[v.ID] = v
	}

	uniq := make(map[string]int)
	var result []dto.TagScanHistoryResponseV2
	for _, v := range mTags {
		tagName := v.Name.String
		if tagName == "" {
			tagName = v.ID
		}
		if resultIdx, ok := uniq[tagName]; !ok {
			uniq[tagName] = len(result)
			result = append(result, dto.TagScanHistoryResponseV2{
				Name:  tagName,
				Count: 1,
			})
		} else {
			result[resultIdx].Count++
		}
	}

	return result
}

func (uc *listGroupByName) buildResponseGroupByDay(mScanHistories []model.TagScanHistory, mTags []model.Tag) []dto.TagScanHistoryResponse {
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

func (uc *listGroupByName) bind(c *gin.Context) (*dto.TagScanHistoryRequest, error) {
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

func (uc *listGroupByName) validate(req *dto.TagScanHistoryRequest) error {
	if err := gvalidator.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}

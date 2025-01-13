package dto

import "time"

type Tag struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type ScanTagRequest struct {
	TagIDs []string `json:"tag_ids" binding:"required,dive"`
}

type TagIDNamePair struct {
	TagID string `json:"tag_id"`
	Name  string `json:"name"`
}

type UpdateTagNameRequest struct {
	TagNames []TagIDNamePair `json:"tag_names" binding:"required,dive"`
}

type TagScanHistoryRequest struct {
	From *time.Time `form:"from" binding:"omitnil,required"`
	To   *time.Time `form:"to" binding:"omitnil,required"`
}

type ScanHistoryUnit struct {
	TagName   string    `json:"tag_name"`
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"created_at"`
}

type TagScanHistoryResponse struct {
	Day       string            `json:"day"`
	Histories []ScanHistoryUnit `json:"histories"`
}

type GetTagByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

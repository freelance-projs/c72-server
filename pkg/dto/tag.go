package dto

import (
	"io"
	"time"
)

type Tag struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	IsScanned bool      `json:"is_scanned"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type ScanTagRequest struct {
	TagIDs []string `json:"tag_ids" validate:"required,dive"`
}

type UpdateTagNameRequest struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type UpdateTagNameRequestByName struct {
	OldName string `json:"old_name" validate:"required"`
	NewName string `json:"new_name" validate:"required"`
}

type GetTagByIDRequest struct {
	ID string `uri:"id" validate:"required"`
}

type ListTagsRequest struct {
	Name      *string    `form:"name" validate:"omitnil,required"`
	IsScanned *bool      `form:"is_scanned", validate:"omitnil,required"`
	From      *time.Time `form:"from" validate:"omitnil,required"`
	To        *time.Time `form:"to" validate:"omitnil,required"`
}

type TagScanHistoryRequest struct {
	From *time.Time `form:"from" validate:"omitnil,required"`
	To   *time.Time `form:"to" validate:"omitnil,required"`
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

type TagScanHistoryResponseV2 struct {
	Name  string `form:"name" json:"name"`
	Count int    `form:"count" json:"count"`
}

type DeleteTagByIDRequest struct {
	ID string `uri:"id" validate:"required"`
}

type DeleteTagByNameRequest struct {
	Name string `uri:"name" validate:"required"`
}

type TagMappingUploadRequest struct {
	Reader io.ReadCloser
}

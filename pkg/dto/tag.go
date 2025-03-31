package dto

import (
	"io"
	"time"
)

type Tag struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type AssignTagRequest struct {
	IDs  []string `json:"ids"`
	Name string   `json:"name"`
}

type ScanTagRequest struct {
	TagIDs []string `json:"tag_ids" validate:"required,dive"`
}

type UpdateTagRequest struct {
	ID   string `json:"id" uri:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type UpdateTagNameRequestByName struct {
	OldName string `json:"old_name" validate:"required"`
	NewName string `json:"new_name" validate:"required"`
}

type GetTagByIDRequest struct {
	ID string `uri:"id" validate:"required"`
}

type ListTagByIDRequest struct {
	IDs  []string `json:"ids" form:"ids" validate:"required,dive"`
	Type string   `json:"type" form:"type"`
}

type ListTagByIDResponse struct {
	Tags            []Tag `json:"tags"`
	DeactivatedTags []Tag `json:"deactivated_tags"`
}

type ListTagsRequest struct {
	Name *string    `form:"name" validate:"omitnil,required"`
	From *time.Time `form:"from" validate:"omitnil,required"`
	To   *time.Time `form:"to" validate:"omitnil,required"`
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

type TagScanHistoryResponseV1 struct {
	Day       string            `json:"day"`
	Histories []ScanHistoryUnit `json:"histories"`
}

type GroupTagByNameResponse struct {
	Name  string `form:"name" json:"name"`
	Count int    `form:"count" json:"count"`
}

type DeleteTagByIDRequest struct {
	ID string `json:"id" uri:"id" validate:"required"`
}

type DeleteTagByNameRequest struct {
	Name string `uri:"name" validate:"required"`
}

type TagMappingUploadRequest struct {
	Reader io.ReadCloser
}

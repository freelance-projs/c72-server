package dto

import "time"

type Lending struct {
	ID          int          `json:"id"`
	Department  string       `json:"department"`
	NumLending  int          `json:"num_lending"`
	NumReturned int          `json:"num_returned"`
	CreatedAt   time.Time    `json:"created_at"`
	Tags        []LendingTag `json:"tags,omitempty"`
}

type DoLendingRequest struct {
	Department string   `json:"department" validate:"required"`
	TagIDs     []string `json:"tag_ids"`
}

type GetLendingRequest struct {
	ID int `json:"id" uri:"id" validate:"required"`
}

type LendingTag struct {
	LendingID int    `json:"lending_id"`
	TagID     string `json:"tag_id"`
	TagName   string `json:"tag_name,omitempty"`
	Status    string `json:"status"`
}

type ReturnDirtyRequest struct {
	TagIDs []string `json:"tag_ids" validate:"required,dive"`
}

type ListLendingRequest struct {
	From      *int64 `json:"from" form:"from" validate:"required,omitnil"`
	To        *int64 `json:"to" form:"to" validate:"required,omitnil"`
	Completed *bool  `json:"completed" form:"completed"`
}

type GetLendingTagsRequest struct {
	LendingID int `json:"lending_id" uri:"id" validate:"required"`
}

type GetWashingTagsRequest struct {
	WashingID int `json:"washing_id" uri:"id" validate:"required"`
}

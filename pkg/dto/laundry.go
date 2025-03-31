package dto

import "time"

type Laundry struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	NumWashing  int          `json:"num_washing"`
	NumReturned int          `json:"num_returned"`
	CreatedAt   time.Time    `json:"created_at"`
	Tags        []LaundryTag `json:"tags,omitempty"`
}

type DoLaundryRequest struct {
	Name   string   `json:"name" validate:"required"`
	TagIDs []string `json:"tag_ids" validate:"required,dive"`
}

type GetLaundryRequest struct {
	ID int `json:"id" uri:"id" validate:"required"`
}

type LaundryTag struct {
	LaundryID int    `json:"laundry_id"`
	TagID     string `json:"tag_id"`
	TagName   string `json:"tag_name,omitempty"`
	Status    string `json:"status"`
}

type ReturnCleanRequest struct {
	TagIDs []string `json:"tag_ids" validate:"required,dive"`
}

type ListLaundryRequest struct {
	From      *int64 `json:"from" form:"from" validate:"required,omitnil"`
	To        *int64 `json:"to" form:"to" validate:"required,omitnil"`
	Completed *bool  `json:"completed" form:"completed"`
}

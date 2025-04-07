package dto

import "time"

type TxLogDepartment struct {
	ID          int           `json:"id"`
	Department  string        `json:"department"`
	NumLending  int           `json:"num_lending"`
	NumReturned int           `json:"num_returned"`
	Details     []TxLogDetail `json:"details,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
}

type ListTxLogDeptRequest struct {
	From *int64 `json:"from" form:"from" validate:"required,omitnil"`
	To   *int64 `json:"to" form:"to" validate:"required,omitnil"`
}

type GetTxLogDeptRequest struct {
	ID int `json:"id" uri:"id" validate:"required"`
}

type TxLogCompany struct {
	ID          int           `json:"id"`
	Company     string        `json:"company"`
	NumWashing  int           `json:"num_washing"`
	NumReturned int           `json:"num_returned"`
	Details     []TxLogDetail `json:"details,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
}

type ListTxLogCompanyRequest struct {
	From *int64 `json:"from" form:"from" validate:"required,omitnil"`
	To   *int64 `json:"to" form:"to" validate:"required,omitnil"`
}

type GetTxLogCompanyRequest struct {
	ID int `json:"id" uri:"id" validate:"required"`
}

type TxLogDetail struct {
	Action    string          `json:"action"`
	Tracking  []TxLogTracking `json:"tracking"`
	CreatedAt time.Time       `json:"created_at"`
}

type TxLogTracking struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type CreateTxLogRequest struct {
	TagIDs []string `json:"tag_ids" validate:"required,dive"`
	Entity string   `json:"entity" validate:"required"`
	Action string   `json:"action" validate:"required,oneof=lending lending_return washing washing_return"`
}

type CreateTxLogResponse struct {
	TxID int `json:"tx_id"`
}

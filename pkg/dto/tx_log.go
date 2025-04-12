package dto

import "time"

type ListTxLogDeptRequest struct {
	From *int64 `json:"from" form:"from" validate:"required,omitnil"`
	To   *int64 `json:"to" form:"to" validate:"required,omitnil"`
}

type CreateTxLogRequest struct {
	TagIDs []string `json:"tag_ids" validate:"required,dive"`
	Entity string   `json:"entity" validate:"required"`
	Action string   `json:"action" validate:"required,oneof=lending lending_return washing washing_return"`
}

type CreateTxLogResponse struct {
	TxID int `json:"tx_id"`
}

type GetTxLogDepartmentRequest struct {
	ID int `json:"id" uri:"id" validate:"required"`
}

type TxLog struct {
	ID       int           `json:"id"`
	Exported int           `json:"exported"`
	Returned int           `json:"returned"`
	Details  []TxLogDetail `json:"details,omitempty"`
}

type TxLogDetail struct {
	Entity    string        `json:"entity"`
	Action    string        `json:"action"`
	Actor     string        `json:"actor"`
	Tracking  []TagTracking `json:"tracking"`
	CreatedAt time.Time     `json:"created_at"`
}

type TagTracking struct {
	Name     string `json:"name"`
	Exported int    `json:"exported"`
	Returned int    `json:"returned"`
}

type ListTxLogRequest struct {
	From *int64 `json:"from" form:"from" validate:"required,omitnil"`
	To   *int64 `json:"to" form:"to" validate:"required,omitnil"`
}

type TxLogDepartment struct {
	ID         int       `json:"id"`
	Department string    `json:"department"`
	Lending    int       `json:"lending"`
	Returned   int       `json:"returned"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetTxLogCompanyRequest struct {
	ID int `json:"id" uri:"id" validate:"required"`
}

type TxLogCompany struct {
	ID        int       `json:"id"`
	Company   string    `json:"company"`
	Washing   int       `json:"washing"`
	Returned  int       `json:"returned"`
	CreatedAt time.Time `json:"created_at"`
}

package model

import "time"

type TxTagStatusEnum string

const (
	TxTagStatusLending TxTagStatusEnum = "lending"
	TxTagStatusWashing TxTagStatusEnum = "washing"
)

type TxTag struct {
	TagID  string          `gorm:"column:tag_id"`
	TxID   int             `gorm:"column:tx_id"`
	Status TxTagStatusEnum `gorm:"column:status"`
}

func (TxTag) TableName() string {
	return "tx_tag"
}

type EDepartmentAction string

const (
	EDepartmentActionLending  EDepartmentAction = "lending"
	EDepartmentActionReturned EDepartmentAction = "returned"
)

func (e EDepartmentAction) String() string {
	return string(e)
}

type TxLogDepartment struct {
	ID         int               `gorm:"column:id"`
	TagName    string            `gorm:"column:tag_name"`
	Department string            `gorm:"column:department"`
	Action     EDepartmentAction `gorm:"column:action"`
	Actor      string            `gorm:"column:actor"`
	Lending    int               `gorm:"column:lending"`
	Returned   int               `gorm:"column:returned"`
	CreatedAt  time.Time         `gorm:"column:created_at"`
}

func (TxLogDepartment) TableName() string {
	return "tx_log_department"
}

type ECompanyAction string

const (
	ECompanyActionWashing  ECompanyAction = "washing"
	ECompanyActionReturned ECompanyAction = "returned"
)

func (e ECompanyAction) String() string {
	return string(e)
}

type TxLogCompany struct {
	ID        int            `gorm:"column:id;primaryKey"`
	Company   string         `gorm:"column:company"`
	Action    ECompanyAction `gorm:"column:action"`
	Actor     string         `gorm:"column:actor"`
	TagName   string         `gorm:"column:tag_name"`
	Washing   int            `gorm:"column:washing"`
	Returned  int            `gorm:"column:returned"`
	CreatedAt time.Time      `gorm:"column:created_at"`
}

func (TxLogCompany) TableName() string {
	return "tx_log_company"
}

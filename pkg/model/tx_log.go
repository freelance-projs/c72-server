package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

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

type TxLogDepartment struct {
	ID        int           `gorm:"primaryKey"`
	Details   TxLogDetails  `gorm:"column:details"`
	Overview  TxLogOverview `gorm:"column:overview"`
	CreatedAt time.Time     `gorm:"column:created_at"`
}

func (TxLogDepartment) TableName() string {
	return "tx_log_dept"
}

type TxLogCompany struct {
	ID        int           `gorm:"primaryKey"`
	Details   TxLogDetails  `gorm:"column:details"`
	Overview  TxLogOverview `gorm:"column:overview"`
	CreatedAt time.Time     `gorm:"column:created_at"`
}

func (TxLogCompany) TableName() string {
	return "tx_log_company"
}

type TxLogDetails []TxLogDetail
type TxLogDetail struct {
	Action    string          `json:"action"`
	Tracking  []TxLogTracking `json:"tracking"`
	CreatedAt time.Time       `json:"created_at"`
}

type TxLogTracking struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (v *TxLogDetails) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to cast TxLogDetail: %v", value)
	}

	result := []TxLogDetail{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return fmt.Errorf("failed to unmarshal TxLogDetail: %w", err)
	}
	*v = result

	return nil
}

func (v TxLogDetails) Value() (driver.Value, error) {
	return json.Marshal(v)
}

type TxLogOverview struct {
	Actor     string `json:"actor"`
	TotalTags uint   `json:"total_tags"`
	Returned  uint   `json:"returned"`
}

func (v *TxLogOverview) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to cast TxLogOverview: %v", value)
	}

	result := TxLogOverview{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return fmt.Errorf("failed to unmarshal TxLogOverview: %w", err)
	}
	*v = result

	return nil
}

func (v TxLogOverview) Value() (driver.Value, error) {
	return json.Marshal(v)
}

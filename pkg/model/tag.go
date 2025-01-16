package model

import (
	"database/sql"
	"time"
)

type Tag struct {
	ID        string         `gorm:"column:id;primaryKey"`
	Name      sql.NullString `gorm:"column:name"`
	IsScanned bool           `gorm:"column:is_scanned"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
}

func (Tag) TableName() string {
	return "tags"
}

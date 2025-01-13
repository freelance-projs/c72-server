package model

import "time"

type TagScanHistory struct {
	ID        int       `gorm:"column:id;primaryKey"`
	TagID     string    `gorm:"column:tag_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (TagScanHistory) TableName() string {
	return "tag_scan_histories"
}

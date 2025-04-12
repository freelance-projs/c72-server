package model

import (
	"database/sql"
)

type Tag struct {
	ID   string         `gorm:"column:id;primaryKey"`
	Name sql.NullString `gorm:"column:name"`
}

func (Tag) TableName() string {
	return "tag"
}

type TagName struct {
	Name string `gorm:"column:name"`
}

func (TagName) TableName() string {
	return "tag_name"
}

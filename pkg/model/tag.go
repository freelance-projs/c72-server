package model

import (
	"database/sql"
	"time"
)

type Tag struct {
	ID        string         `gorm:"column:id;primaryKey"`
	Name      sql.NullString `gorm:"column:name"`
	CreatedAt time.Time      `gorm:"column:created_at"`
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

type TagDepartment struct {
	ID   string `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name"`
}

func (TagDepartment) TableName() string {
	return "tag_department"
}

type TagCompany struct {
	ID   string `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name"`
}

func (TagCompany) TableName() string {
	return "tag_company"
}

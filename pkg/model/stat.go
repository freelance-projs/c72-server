package model

import "time"

type LendingStat struct {
	ID         int       `gorm:"column:id;primaryKey"`
	TagName    string    `gorm:"column:tag_name;primaryKey"`
	Department string    `gorm:"column:department"`
	Lending    int       `gorm:"column:lending"`
	Returned   int       `gorm:"column:returned"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (LendingStat) TableName() string {
	return "lending_stat"
}

type WashingStat struct {
	ID        int       `gorm:"column:id;primaryKey"`
	TagName   string    `gorm:"column:tag_name;primaryKey"`
	Company   string    `gorm:"column:company"`
	Washing   int       `gorm:"column:washing"`
	Returned  int       `gorm:"column:returned"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (WashingStat) TableName() string {
	return "washing_stat"
}

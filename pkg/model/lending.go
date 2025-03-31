package model

import "time"

type Lending struct {
	ID         int       `gorm:"column:id;primaryKey"`
	Department string    `gorm:"column:department"`
	NumLending int       `gorm:"num_lending"`
	CreatedAt  time.Time `gorm:"column:created_at"`

	Tags []LendingTag `gorm:"foreignKey:LendingID;references:ID"`
}

func (Lending) TableName() string {
	return "lending"
}

func (Lending) TagsRelation() string {
	return "Tags"
}

type lendingColumns struct {
	ID         string
	Department string
	NumLending string
	CreatedAt  string
}

func (Lending) Columns() lendingColumns {
	return lendingColumns{
		ID:         "id",
		Department: "department",
		NumLending: "num_lending",
		CreatedAt:  "created_at",
	}
}

type LendingTagStatusEnum string

const (
	LendingStatusLending  LendingTagStatusEnum = "lending"
	LendingStatusReturned LendingTagStatusEnum = "returned"
)

func (e LendingTagStatusEnum) String() string {
	return string(e)
}

// todo: partition by week
type LendingTag struct {
	LendingID int                  `gorm:"column:lending_id"`
	TagID     string               `gorm:"column:tag_id"`
	Status    LendingTagStatusEnum `gorm:"column:status"`
	Tag       *Tag                 `gorm:"foreignKey:TagID;references:ID"`
	CreatedAt time.Time            `gorm:"column:created_at"`
}

func (LendingTag) TableName() string {
	return "lending_tag"
}

func (LendingTag) TagRelation() string {
	return "Tag"
}

type lendingTagColumns struct {
	LendingID string
	TagID     string
	Status    string
}

func (LendingTag) Columns() lendingTagColumns {
	return lendingTagColumns{
		LendingID: "lending_id",
		TagID:     "tag_id",
		Status:    "status",
	}
}

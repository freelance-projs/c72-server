package model

import "time"

type Laundry struct {
	ID         int       `gorm:"column:id;primaryKey"`
	Name       string    `gorm:"column:name"`
	NumWashing int       `gorm:"column:num_washing"`
	CreatedAt  time.Time `gorm:"column:created_at"`

	Tags []LaundryTag `gorm:"foreignKey:LaundryID;references:ID"`
}

func (Laundry) TableName() string {
	return "laundry"
}

func (Laundry) TagsRelation() string {
	return "Tags"
}

type laundryColumns struct {
	ID         string
	Name       string
	NumWashing string
	CreatedAt  string
}

func (Laundry) Columns() laundryColumns {
	return laundryColumns{
		ID:         "id",
		Name:       "name",
		NumWashing: "num_washing",
		CreatedAt:  "created_at",
	}
}

type LaundryTag struct {
	LaundryID int               `gorm:"column:laundry_id"`
	TagID     string            `gorm:"column:tag_id"`
	Status    LaundryStatusEnum `gorm:"column:status"`

	Tag *Tag `gorm:"foreignKey:TagID;references:ID"`
}

func (LaundryTag) TableName() string {
	return "laundry_tag"
}

func (LaundryTag) TagRelation() string {
	return "Tag"
}

type laundryTagColumns struct {
	LaundryID string
	TagID     string
	Status    string
}

func (LaundryTag) Columns() laundryTagColumns {
	return laundryTagColumns{
		LaundryID: "laundry_id",
		TagID:     "tag_id",
		Status:    "status",
	}
}

type LaundryStatusEnum string

const (
	LaundryStatusWashing  LaundryStatusEnum = "washing"
	LaundryStatusReturned LaundryStatusEnum = "returned"
)

func (e LaundryStatusEnum) String() string {
	return string(e)
}

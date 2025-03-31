package repository

import (
	"gorm.io/gorm"
)

type Laundry struct {
	db *gorm.DB
}

func NewLaundry(db *gorm.DB) *Laundry {
	return &Laundry{db: db}
}

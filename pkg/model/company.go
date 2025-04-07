package model

type Company struct {
	Name string `gorm:"column:name;primaryKey"`
}

func (Company) TableName() string {
	return "company"
}

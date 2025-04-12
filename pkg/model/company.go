package model

type Company struct {
	Name string `gorm:"column:name;primaryKey"`
}

func (Company) TableName() string {
	return "company"
}

type CompanyTag struct {
	TagID   string `gorm:"column:tag_id;primaryKey"`
	Company string `gorm:"column:company"`
}

func (CompanyTag) TableName() string {
	return "company_tag"
}

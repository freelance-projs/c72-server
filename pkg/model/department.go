package model

type Department struct {
	Name string `gorm:"column:name;primaryKey"`
}

func (Department) TableName() string {
	return "department"
}

type DepartmentTag struct {
	TagID      string `gorm:"column:tag_id;primaryKey"`
	Department string `gorm:"column:department"`
}

func (DepartmentTag) TableName() string {
	return "department_tag"
}

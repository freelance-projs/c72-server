package model

type Department struct {
	Name string `gorm:"column:name;primaryKey"`
}

func (Department) TableName() string {
	return "department"
}

func (Department) Columns() departmentColumns {
	return departmentColumns{
		Name: "name",
	}
}

type departmentColumns struct {
	Name string
}

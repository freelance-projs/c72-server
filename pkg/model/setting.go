package model

type Setting struct {
	Key   string `gorm:"column:key;primaryKey"`
	Value string `gorm:"column:value"`
}

func (s Setting) TableName() string {
	return "setting"
}

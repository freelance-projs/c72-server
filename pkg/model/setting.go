package model

type Setting struct {
	ID            int    `gorm:"column:id;primaryKey"`
	TxLogSheetID  string `gorm:"column:tx_log_sheet_id"`
	ReportSheetID string `gorm:"column:report_sheet_id"`
}

func (s Setting) TableName() string {
	return "setting"
}

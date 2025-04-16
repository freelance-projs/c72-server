package dto

import "io"

type TagName struct {
	Name string `json:"name"`
}

type CreateBatchTagNameRequest struct {
	Reader io.ReadCloser
}

type ListTagNameRequest struct{}

type UpdateSettingRequest struct {
	TxLogSheetID  string `json:"tx_log_sheet_id" validate:"omitempty,min=1"`
	ReportSheetID string `json:"report_sheet_id" validate:"omitempty,min=1"`
}

type DeleteTagNameRequest struct {
	Names []string `json:"names" validate:"required,dive"`
}

type ChangeTagNameRequest struct {
	OldName string `json:"old_name" validate:"required"`
	NewName string `json:"new_name" validate:"required"`
}

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
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type DeleteTagNameRequest struct {
	Names []string `json:"names" validate:"required,dive"`
}

type ChangeTagNameRequest struct {
	OldName string `json:"old_name" validate:"required"`
	NewName string `json:"new_name" validate:"required"`
}

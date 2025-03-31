package dto

import "io"

type Department struct {
	Name string `json:"name"`
}

type CreateDepartmentRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateBatchDepartmentRequest struct {
	Reader io.ReadCloser
}

type ListDepartmentsRequest struct{}

type DeleteDepartmentRequest struct {
	Names []string `json:"names" validate:"required,dive"`
}

type ChangeDepartmentRequest struct {
	OldName string `json:"old_name" validate:"required"`
	NewName string `json:"new_name" validate:"required"`
}

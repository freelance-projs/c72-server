package dto

import "io"

type Company struct {
	Name string `json:"name"`
}

type CreateCompanyRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateBatchCompanyRequest struct {
	Reader io.ReadCloser
}

type ListCompaniesRequest struct {
}

type DeleteCompanyRequest struct {
	Names []string `json:"names" validate:"required,dive"`
}

type ChangeCompanyRequest struct {
	OldName string `json:"old_name" validate:"required"`
	NewName string `json:"new_name" validate:"required"`
}

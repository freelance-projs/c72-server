package dto

type Setting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateSettingRequest struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type DeleteSettingRequest struct {
	Key string `json:"key" validate:"required"`
}

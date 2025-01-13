package dto

type LogRequest struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

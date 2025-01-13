package dto

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/common/apperror"
)

// Response represent common format response to client
type Response struct {
	Success    bool   `json:"success"`
	Data       any    `json:"data,omitempty"`
	StatusCode int    `json:"-"`
	Message    string `json:"message,omitempty"`
	Error      any    `json:"error,omitempty"`
	Paging     any    `json:"paging,omitempty"`
}

func StatusOK(data any, opts ...func(*Response)) *Response {
	res := &Response{
		Data:       data,
		Success:    true,
		StatusCode: http.StatusOK,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(res)
		}
	}

	return res
}

func StatusCreated(data any, entity string) *Response {
	return &Response{
		Data:       data,
		Success:    true,
		StatusCode: http.StatusCreated,
		Message:    fmt.Sprintf("%s is created", entity),
	}
}

func JSONSuccess(c *gin.Context, respDTO *Response) {
	if respDTO.StatusCode == 0 {
		respDTO.StatusCode = http.StatusOK
	}

	c.JSON(respDTO.StatusCode, respDTO)
}

func JSONFail(c *gin.Context, err error) {
	if baseErr, ok := asHTTPErr(c, err); ok {
		logBaseErr(baseErr)
		return
	}
	if baseErr, ok := asBaseErr(c, err); ok {
		logBaseErr(baseErr)
		return
	}

	httpError := apperror.ErrInternalServer(err)
	logBaseErr(&httpError.BaseError)

	c.JSON(httpError.HTTPCode, Response{
		Success: false,
		Message: httpError.Error(),
	})
}

func logBaseErr(baseErr *apperror.BaseError) {
	kvs := []any{"err_id", baseErr.ID}
	if baseErr.Ancestor() != nil {
		kvs = append(kvs, "err", baseErr.Ancestor())
	}
	slog.Error(baseErr.Error(), kvs...)
}

func JSONAbort(c *gin.Context, err error) {
	c.Abort()
	JSONFail(c, err)
}

func asHTTPErr(c *gin.Context, err error) (*apperror.BaseError, bool) {
	var httpErr *apperror.HTTPError
	if errors.As(err, &httpErr) {

		c.JSON(httpErr.HTTPCode, Response{
			Success: false,
			Error:   err,
			Message: httpErr.Error(),
		})
		return &httpErr.BaseError, true
	}

	return nil, false
}

func asBaseErr(c *gin.Context, err error) (*apperror.BaseError, bool) {
	var baseErr *apperror.BaseError
	if errors.As(err, &baseErr) {

		httpErr := apperror.NewHTTPError(err, http.StatusBadRequest)
		c.JSON(httpErr.HTTPCode, Response{
			Success: false,
			Error:   err,
			Message: baseErr.Error(),
		})
		return baseErr, true
	}

	return nil, false
}

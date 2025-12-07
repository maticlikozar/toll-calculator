package errors

import (
	"fmt"
	"net/http"
)

func NewAPIError(code int, message string) *APIError {
	return &APIError{
		StatusCode: code,
		Title:      message,
		ErrorCode:  code,
	}
}

// APIError represents API errors.
type APIError struct {
	Title      string
	Detail     string
	StatusCode int
	ErrorCode  int
}

// WithDetails set title of the error.
func (a *APIError) WithTitle(title string) *APIError {
	return &APIError{
		StatusCode: a.StatusCode,
		Title:      title,
		Detail:     a.Detail,
		ErrorCode:  a.ErrorCode,
	}
}

// WithDetails set details of the error.
func (a *APIError) WithDetails(details string, args ...interface{}) *APIError {
	return &APIError{
		StatusCode: a.StatusCode,
		Title:      a.Title,
		Detail:     fmt.Sprintf(details, args...),
		ErrorCode:  a.ErrorCode,
	}
}

// WithErrorCode set error code of the error.
func (a *APIError) WithErrorCode(errorCode int) *APIError {
	return &APIError{
		StatusCode: a.StatusCode,
		Title:      a.Title,
		Detail:     a.Detail,
		ErrorCode:  errorCode,
	}
}

// Error returns string representation of the error.
func (e *APIError) Error() string {
	return fmt.Sprintf("ErrorCode: %d: Title: %s, Detail: %s, ErrorCode: %d", e.ErrorCode, e.Title, e.Detail, e.ErrorCode)
}

var (
	ErrAPIUnauthorized = NewAPIError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	ErrAPIForbidden    = NewAPIError(http.StatusForbidden, http.StatusText(http.StatusForbidden))
	ErrAPIBadRequest   = NewAPIError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	ErrAPINotFound     = NewAPIError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	ErrAPIInternal     = NewAPIError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
)

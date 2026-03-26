package apperror

import "fmt"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	HTTP    int    `json:"-"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("[%d] %s", e.Code, e.Message)
	}
	return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
}

func (e *AppError) Unwrap() error { return e.Err }

var (
	ErrBadRequest   = &AppError{Code: 4000, Message: "bad request", HTTP: 400}
	ErrUnauthorized = &AppError{Code: 4001, Message: "unauthorized", HTTP: 401}
	ErrForbidden    = &AppError{Code: 4003, Message: "forbidden", HTTP: 403}
	ErrNotFound     = &AppError{Code: 4004, Message: "not found", HTTP: 404}
	ErrConflict     = &AppError{Code: 4009, Message: "conflict", HTTP: 409}
	ErrInternal     = &AppError{Code: 5000, Message: "internal server error", HTTP: 500}
)

func Wrap(base *AppError, err error, msg string) *AppError {
	appErr := *base
	appErr.Err = err
	if msg != "" {
		appErr.Message = msg
	}
	return &appErr
}

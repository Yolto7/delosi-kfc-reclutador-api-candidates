package error

import (
	"fmt"
	"net/http"
)

// ==== TYPES ====
type ErrorType string

const (
	BAD_REQUEST          ErrorType = "BAD_REQUEST"
	UNAUTHORIZED         ErrorType = "UNAUTHORIZED"
	FORBIDDEN            ErrorType = "FORBIDDEN"
	NOT_FOUND            ErrorType = "NOT_FOUND"
	NOT_ALLOWED          ErrorType = "NOT_ALLOWED"
	UNSUPPORTED_CONTENT  ErrorType = "UNSUPPORTED_CONTENT_TYPE"
	UNPROCESSABLE_ENTITY ErrorType = "UNPROCESSABLE_ENTITY"
)

var errorTypeToHttpCode = map[ErrorType]int{
	BAD_REQUEST:          http.StatusBadRequest,
	UNAUTHORIZED:         http.StatusUnauthorized,
	FORBIDDEN:            http.StatusForbidden,
	NOT_FOUND:            http.StatusNotFound,
	NOT_ALLOWED:          http.StatusMethodNotAllowed,
	UNSUPPORTED_CONTENT:  http.StatusUnsupportedMediaType,
	UNPROCESSABLE_ENTITY: http.StatusUnprocessableEntity,
}

// ==== AppError ====
type AppError struct {
	HttpCode  int       `json:"-"`
	ErrorType ErrorType `json:"errorType"`
	Message   string    `json:"message"`
	ErrorCode string    `json:"errorCode,omitempty"`
	Payload   any       `json:"payload,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Creates an AppError directly (for controllers or services)
func New(errorType ErrorType, message, errorCode string, payload ...any) *AppError {
	code, ok := errorTypeToHttpCode[errorType]
	if !ok {
		code = http.StatusInternalServerError
	}

	var p any
	if len(payload) > 0 {
		p = payload[0]
	}

	return &AppError{
		HttpCode:  code,
		ErrorType: errorType,
		Message:   message,
		ErrorCode: errorCode,
		Payload:   p,
	}
}

// ==== CustomError: Reusable error for multiple services ====
type CustomError struct {
	ErrorType ErrorType
	Message   string
	ErrorCode string
	Payload   any
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode, e.Message)
}

// Constructor
func NewError(errorType ErrorType, message, errorCode string, payload ...any) *CustomError {
	var p any
	if len(payload) > 0 {
		p = payload[0]
	}

	return &CustomError{
		ErrorType: errorType,
		Message:   message,
		ErrorCode: errorCode,
		Payload:   p,
	}
}

// ==== Generic error transformer -> AppError ====
func FromError(err error) *AppError {
	if err == nil {
		return nil
	}

	// If is a CustomError, convert it to AppError
	if cerr, ok := err.(*CustomError); ok {
		return New(cerr.ErrorType, cerr.Message, cerr.ErrorCode, cerr.Payload)
	}

	// If is already an AppError, return it as is
	if aerr, ok := err.(*AppError); ok {
		return aerr
	}

	// Generic error
	return New(BAD_REQUEST, err.Error(), "GENERIC_ERROR")
}


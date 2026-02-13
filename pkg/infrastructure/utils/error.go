package utils

import (
	"runtime/debug"
)

type SafeError struct {
	Message string `json:"message"`
	Error   error `json:"error"`
	Stack   string `json:"stack"`
}

func NewSafeError(err error, msg string) SafeError {
	return SafeError{
		Message: msg,
		Error:   err,
		Stack:   string(debug.Stack()),
	}
}
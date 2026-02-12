package rest

import "fmt"

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func CommonError(err error, message string) *Error {
	return &Error{
		Code:    "error",
		Message: fmt.Sprintf("%s: %s", message, err.Error()),
	}
}

func ErrValidate(err error) *Error {
	return &Error{
		Code:    "validation error",
		Message: err.Error(),
	}
}

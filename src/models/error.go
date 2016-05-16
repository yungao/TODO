package models

import (
	"fmt"
)

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (err *Error) Error() string {
	if err == nil {
		return ""
	}

	return fmt.Sprintf("[%d] %s", err.Code, err.Msg)
}

// NewError creates an error instance with the specified code and message.
func NewError(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

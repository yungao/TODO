package models

import (
	"fmt"
)

const (
	// request failed
    ERR_REQUEST_FAILED = 40000
    // invalid data
	ERR_INVALID_DATA   = 40001
    // user disabled
    ERR_USER_DISABLED  = 40002

    // ======== User ==========
	ERR_NAME_EXIST     = 40011
	ERR_USER_NOT_FOUND = 40012
    // ========================

    // ======== Todo ==========
	ERR_TODO_NOT_FOUND = 40111
	ERR_TODO_DISABLED  = 40112
	ERR_TODO_USER_PERM = 40113
    // ========================

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

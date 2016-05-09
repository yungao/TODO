package model

import (
    "fmt"
)

type Error struct {
    Errno   int             `json:"errno"`
    Msg     string          `json:"msg"`
}

func (err *Error) Error() string {
    if err == nil {
        return ""
    }

    return fmt.Sprintf("errno:%d, msg:%s", err.Errno, err.Msg)
}

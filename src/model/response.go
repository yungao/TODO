package model

import (
    "fmt"
)

type Response struct {
    Result  bool            `json:"result"`
    Errno   int             `json:"errno"`
    Msg     string          `json:"msg"`
    Data    interface{}     `json:"data"`
}

func (rp *Response) Response() string {
    if rp == nil {
        return ""
    }

    result := "false"
    if rp.Result {
        result = "true"
    }

    if rp.Data == nil {
        return fmt.Sprintf("result:%s, errno:%d, msg:%s, data: %s", result, rp.Errno, rp.Msg, "")
    } else {
        return fmt.Sprintf("result:%s, errno:%d, msg:%s, data: %s", result, rp.Errno, rp.Msg, rp.Data)
    }
}

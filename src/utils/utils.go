package utils

import (
    "log"
    "strconv"
    "strings"

    "github.com/martini-contrib/sessions"
    "github.com/martini-contrib/render"
)

func ParseSession(session sessions.Session, render render.Render) int64 {
    _id := session.Get("ID")

    id, ok := _id.(int64)
    if !ok {
        return id
    }

    idStr, ok := _id.(string)
    if !ok {
        id, err := strconv.ParseInt(idStr, 0, 64)
        if err == nil {
            return id
        }
    }

    log.Println("Parse session error, unthenticated!")
    render.JSON(401, "Unauthorized!")
    return -1
}

func IsEmpty(s string) bool {
    return (strings.TrimSpace(s) == "")
}

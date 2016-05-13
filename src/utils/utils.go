package utils

import (
    "log"
    "strconv"
    "strings"
    "errors"

    "github.com/martini-contrib/sessions"
    "github.com/martini-contrib/render"
)

func ParseSession(session sessions.Session, render render.Render) (int64, error) {
    if session != nil {
        sid := session.Get("ID")
        log.Printf("Session, ID: %s", sid)

        if id, ok := sid.(int64); ok {
            log.Printf("Parse session, ID: %d", id)
            return id, nil
        }

        if id, ok := sid.(string); ok {
            id, err := strconv.ParseInt(id, 0, 64)
            if err == nil {
                log.Printf(">Parse session, ID: %d", id)
                return id, nil
            }
        }

    }

    log.Println("Parse session error, unthenticated!")
    render.JSON(401, "Unauthorized")

    return -1, errors.New("Unauthorized")
}

func IsEmpty(s string) bool {
    return (strings.TrimSpace(s) == "")
}

package utils

import (
	"encoding/base64"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
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

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64Decode(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

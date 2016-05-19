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

func ParseSession(session sessions.Session, render render.Render) (int, string, error) {
	if session != nil {
		value := session.Get("ID")
		if v, ok := value.(string); ok {
			s := strings.Split(v, ":")
			if len(s) == 2 {
				sid := s[0]
				name := s[1]
				log.Printf("Session, Name: %s[%s]", name, sid)

				id, err := strconv.Atoi(sid)
				if err == nil {
					log.Println("Parse session, ID: ", id)
					return id, name, nil
				}
			}
		}
	}

	log.Println("Parse session error, unthenticated!")
	render.JSON(401, "Unauthorized")

	return -1, "", errors.New("Unauthorized")
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

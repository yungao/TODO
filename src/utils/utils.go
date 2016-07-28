package utils

import (
	"encoding/base64"
	"errors"
	"log"
	"strconv"
	"strings"
	"net/http"

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
				log.Printf("Session -> Name: %s[%s]", name, sid)

				id, err := strconv.Atoi(sid)
				if err == nil {
					return id, name, nil
				}
			}
		}
	}

	log.Println("Parse session error, unthenticated!")
	render.JSON(401, "Unauthorized")

	return -1, "", errors.New("Unauthorized")
}

func ParseUserAgent(request *http.Request) string {
    agent := strings.ToLower(request.Header.Get("User-Agent"))
    if agent == "" {
        return "UFO"
    } else {
        if strings.Contains(agent, "android") {
            return "Android"
        } else if strings.Contains(agent, "iphone") {
            return "iPhone"
        } else if strings.Contains(agent, "mobile") {
            return "Mobile"
        } else if strings.Contains(agent, "msie") {
            return "MSIE"
        } else if strings.Contains(agent, "chrome") {
            return "Chrome"
        } else if strings.Contains(agent, "firefox") {
            return "FireFox"
        } else if strings.Contains(agent, "trident") && strings.Contains(agent, "like gecko") {
            return "MSIE 11"
        } else {
            return "Browser"
        }
    }
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

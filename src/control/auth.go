package control

import (
    "os"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/auth"
)

func Auth() martini.Handler {
    return auth.Basic(os.Getenv("API_USERNAME"), os.Getenv("API_PASSWORD"))
}

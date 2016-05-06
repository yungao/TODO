package main

import (
    "fmt"
    //"log"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/binding"
    "github.com/martini-contrib/render"
    //"github.com/codegangsta/envy/lib"

    "model"
    "config"
    "control"
)

func main() {
    app := martini.Classic()
    // initialization database
    app.Map(config.DB())
    // set asset directory
    app.Use(martini.Static("assets"))

    app.Use(render.Renderer())

    //=================== TODO Router ======================
    app.Get("/", func() string {
        return "Hi, This is TODO Go Web Server!"
    })
    app.Group("/user", func(router martini.Router) {
        router.Get("/:id", control.GetUser)
        router.Post("", binding.Bind(model.User{}), control.CreateUser)
    })//, control.Auth())
    //======================================================

    fmt.Printf("Go web server running...\n")
    app.Run()
    fmt.Printf("Go web server stop !\n")
}

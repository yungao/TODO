package main

import (
    "log"
    "net/http"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/binding"
    "github.com/martini-contrib/render"
    "github.com/martini-contrib/sessions"
    //"github.com/codegangsta/envy/lib"

    "model"
    "config"
    "controller"
)

func main() {
    app := martini.Classic()
    // initialization database
    app.Map(config.DB())
    // set asset directory
    app.Use(martini.Static("assets"))
    // use martini-contrib/render
    app.Use(render.Renderer())
    // use martini-contrib/sessions
    store := sessions.NewCookieStore([]byte("todo2016@eli"))
    app.Use(sessions.Sessions("todo_session", store))

    // 验证session
    app.Use(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.Path != "/api/v1/user/login" {
            if req.Header.Get("Cookies") != "gaoyun" {
                res.WriteHeader(http.StatusUnauthorized)
            } else {
                req.Header.Set("User", "gaoyun")
            }
        }
    })

    //=================== TODO Router ======================
    app.Get("/", func() string {
        return "Hi, This is TODO Go Web Server!"
    })

    app.Group("/api/v1", func(router martini.Router) {
        app.Group("/user", func(router martini.Router) {
            router.Post("/login", func() string {
                return "Welcome login TODO!"
            })
            router.Get("", controller.ListUser)
            router.Get("/:id", controller.GetUser)
            router.Post("", binding.Bind(model.User{}), controller.CreateUser)
            router.Delete("/:id", controller.DeleteUser)
            router.Patch("/pwd/:id", controller.UpdateUser)
        })
    })
    //======================================================

    log.Print("Go web server running...\n")
    app.Run()
    log.Print("Go web server stop !\n")
}

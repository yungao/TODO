package main

import (
	"log"
    "os"
	//"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
    "github.com/coopernurse/gorp"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/sessionauth"
	//"github.com/codegangsta/envy/lib"

	config      "config"
	control     "controllers"
	model       "models"
)


func initDBTables(db *gorp.DbMap) {
    db.AddTableWithName(model.User{}, "user").SetKeys(true, "ID")
    db.CreateTables()
    //db.DropTables()
}

func enableDBLogger(db *gorp.DbMap) {
    db.TraceOn("[gorp]", log.New(os.Stdout, ">>MySQL<< ", log.Lmicroseconds))
    //dbp.TraceOff()
}


func main() {
    //envy.Bootstrap()

	app := martini.Classic()
	// initialization database
    dbmap := config.DB()
	app.Map(dbmap)
    initDBTables(dbmap);
    enableDBLogger(dbmap);

    // set asset directory
    app.Use(martini.Static("assets"))
    // use martini-contrib/render
    app.Use(render.Renderer())
    // use martini-contrib/sessions
    store := sessions.NewCookieStore([]byte("secret123"))
    // Default our store to use Session cookies, so we don't leave logged in
    // users roaming around
    store.Options(sessions.Options{
        MaxAge: 0,
    })
    app.Use(sessions.Sessions("my_session", store))
    app.Use(sessionauth.SessionUser(model.GenerateAnonymousUser))
    sessionauth.RedirectUrl = "/api/v1/unauth"
    sessionauth.RedirectParam = "redirect"

	//=================== TODO Router ======================
	app.Get("/", func(render render.Render) string {
        // render.HTML(200, "index", nil)
		return "Hi, This is TODO Go Web Server!"
	})

    app.Group("/api/v1", func(router martini.Router) {
        app.Get("/unauth", func(render render.Render) {
            // render.HTML(200, "login", nil)
            render.JSON(401, "Please Login TODO first!")
        })

        app.Group("/user", func(router martini.Router) {
            router.Post("/login", binding.Bind(model.User{}), control.Login)
            router.Get("", control.ListUser)
            router.Get("/:id([(0-9)+])", control.GetUser)
            router.Post("", binding.Bind(model.User{}), control.CreateUser)
            router.Delete("/:id", control.DeleteUser)
            router.Patch("/pwd/:id", control.UpdateUser)
            router.Get("/isLogin", sessionauth.LoginRequired, func(render render.Render, user sessionauth.User) {
                render.JSON(200, user)
            })
        })
	})
	//======================================================

	log.Print("Go web server running...\n")
	app.Run()
	log.Print("Go web server stopped !\n")
}

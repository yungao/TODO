package main

import (
	"log"
	// "os"
	// "net/http"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	// "github.com/codegangsta/envy/lib"

	config "config"
	control "controllers"
	model "models"
)

func initDBTables(db *gorp.DbMap) {
	model.CreateUserTable(db)
	model.CreatePartnerTable(db)
	model.CreateTodoTable(db)
	model.CreateTagTable(db)
	model.CreateTodoTagTable(db)
	model.CreateProcessTable(db)
	model.CreateAttachTable(db)
}

func main() {
	//envy.Bootstrap()

	app := martini.Classic()
	// initialization database
	db := config.DB()
	initDBTables(db)
	config.EnableDBLogger(db)
	app.Map(db)

	// set asset directory
	app.Use(martini.Static("assets"))
	// use martini logger
	app.Use(martini.Logger())
	// use martini-contrib/render
	app.Use(render.Renderer())
	// use martini-contrib/sessions
	store := sessions.NewCookieStore([]byte("todo2016@etech"))
	// Default our store to use Session cookies, so we don't leave logged in
	// users roaming around
	store.Options(sessions.Options{
		MaxAge: 0,
	})
	app.Use(sessions.Sessions("session", store))

	//====================== Router ========================
	app.Get("/", func(render render.Render) string {
		// render.HTML(200, "index", nil)
		return "Hi, This is TODO Go Web Server!"
	})

	app.Group("/api/v1", func(router martini.Router) {
		// Http request unauthorized
		// app.Get("/unauth", func(render render.Render) {
		//    // render.HTML(200, "login", nil)
		//    render.JSON(401, "Please Login TODO first!")
		// })

		app.Group("/user", func(router martini.Router) {
			// registered user
			router.Post("", binding.Bind(model.User{}), control.CreateUser)
			// update user info
			router.Patch("", control.UpdateUser)
			router.Patch("/:id([(0-9)+])", control.UpdateUser)

			// List user informations
			router.Get("", control.ListUser)

			// User login
			router.Post("/login", binding.Bind(model.User{}), control.Login)
			router.Get("/login", binding.Bind(model.User{}), control.Login)
			// Get user information
			router.Get("/:id([(0-9)+])", control.GetUser)
			// Delete user
			router.Delete("/:id", control.DeleteUser)
			// Update user
			router.Patch("/pwd/:id", control.UpdateUser)
			// Check login state
			router.Get("/isLogin", func(session sessions.Session, render render.Render) {
				v := session.Get("ID")
				render.JSON(200, v)
			})
		})
	})
	app.NotFound(func(render render.Render) {
		render.JSON(404, "Not Found!")
	})
	//======================================================

	log.Print("Go web server running...\n")
	app.Run()
	log.Print("Go web server stopped !\n")
}

package main

import (
	"log"
	"os"
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

func genHttpPort() string {
	switch os.Getenv("HOME") {
	case "/home/gaoyun":
		return ":3000"
	case "/home/lquan":
		return ":3001"
	case "/home/yxd":
		return ":3002"
	default:
		return ":3010"
	}
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
	app.Use(martini.Static("gui"))
	// use martini logger
	// app.Use(martini.Logger())
	// use martini-contrib/render
	app.Use(render.Renderer())
	// use martini-contrib/sessions
	store := sessions.NewCookieStore([]byte("todo2016@etech"))
	store.Options(sessions.Options{
		MaxAge: 2 * 60 * 60, // 2*60*60s
        Path: "/",
	})
	app.Use(sessions.Sessions("session", store))

	//====================== Router ========================
	app.Get("/", func(render render.Render) string {
		// render.HTML(200, "index", nil)
		return "Hi, This is TODO Go Web Server!"
	})

	app.Group("/api/v1", func(router martini.Router) {
		app.Group("/user", func(router martini.Router) {
			// registered user
			router.Post("", binding.Bind(model.User{}), control.CreateUser)
			// update user info
			router.Patch("/:id", control.UpdateUser)
			router.Patch("", control.UpdateUser)
			// user login
			router.Post("/login", binding.Bind(model.User{}), control.Login)
			// router.Get("/login", binding.Bind(model.User{}), control.Login)
			// check login state
			router.Get("/login", control.IsLogin)
			// user logout
			router.Get("/logout", control.Logout)
			// get user info
			router.Get("/:id", control.GetUser)
			// list users info, can filter
			router.Get("", control.ListUsers)
			// delete user
			router.Delete("/:id", control.DeleteUser)
		})

		app.Group("/todo", func(router martini.Router) {
			// create todo
			router.Post("", binding.Bind(model.Todo{}), control.CreateTodo)
			// get todo info
			router.Get("/:id", control.GetTodo)
			// list todo, can filter
			router.Get("", control.ListTodos)
			// update todo
			router.Patch("/:id", control.UpdateTodo)

			// add todo partner
			router.Post("/partner", binding.Bind(model.Partner{}), control.AddPartner)
			router.Post("/tag", binding.Bind(model.TodoTag{}), control.AddTodoTag)
		})

		app.Group("/tag", func(router martini.Router) {
            // create tag
			router.Post("", binding.Bind(model.Tag{}), control.CreateTag)
		})
	})
	app.NotFound(func(render render.Render) {
		render.JSON(404, "Not Found!")
	})
	//======================================================

	log.Print("Go web server running...\n")
	// app.Run()
	app.RunOnAddr(genHttpPort())
	log.Print("Go web server stopped !\n")
}

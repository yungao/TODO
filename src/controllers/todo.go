package controllers

import (
	"fmt"
    "time"
	"log"
	"net/http"
	 "strconv"
	// "strings"
	// "regexp"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	//"github.com/martini-contrib/binding"

	// "config"
	model "models"
	utils "utils"
)

/**
* Create a new Todo
 */
func CreateTodo(session sessions.Session, todo model.Todo, db *gorp.DbMap, render render.Render) {
	log.Println("Create Todo: ", todo.String())
	id, _, err := utils.ParseSession(session, render)
    if err == nil {
        user, err := model.GetUserByID(db, id)
        if err == nil && user.IsEnable() {
            if len(todo.Name) > 100 {
                render.JSON(422, model.NewError(model.ERR_INVALID_DATA, "Todo name must be 1-100 characters!"))
                return
            }
            if len(todo.Content) > 2048 {
                render.JSON(422, model.NewError(model.ERR_INVALID_DATA, "Todo content must be 0-2048 characters!"))
                return
            }
            if todo.LimitAt <= time.Now().Unix() {
                render.JSON(422, model.NewError(model.ERR_INVALID_DATA, "Invalid limit time!"))
                return
            }

            todo.CreatorID = id

            err = db.Insert(&todo)
            if err == nil {
                // TODO... add todo process
                render.JSON(201, &todo)
            } else {
                log.Printf("Create todo error: %s", err.Error())
                erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: err.Error()}
                render.JSON(422, erp)
            }
        } else {
            erp := model.Error{Code: model.ERR_USER_DISABLED, Msg: "User is disabled!"}
            render.JSON(403, erp)
        }
    }
}

/**
* Get todo info
 */
func GetTodo(session sessions.Session, db *gorp.DbMap, params martini.Params, render render.Render) {
	_, _, err := utils.ParseSession(session, render)
	if err == nil { // has login
		param := params["id"]
		id, err := strconv.Atoi(param)
		if err == nil { // get by ID
			todo, err := model.GetTodoByID(db, id)
			if err == nil {
				render.JSON(200, todo)
				return
			}
		}

		render.JSON(404, model.NewError(model.ERR_TODO_NOT_FOUND, fmt.Sprintf("Todo[%s] does not exist!", param)))
	}
}

/**
* List todo
 */
func ListTodos(session sessions.Session, db *gorp.DbMap, params martini.Params, render render.Render, request *http.Request) {
	_, _, err := utils.ParseSession(session, render)
	if err == nil {         // has login
        todos, err := model.GetAllTodos(db)
        if err == nil {
            render.JSON(200, todos)
        }

		//query := request.URL.Query()
		//if len(query) == 0 {    // list all todos
		//	var users, err = model.GetAllUsers(db)
		//	if err != nil {
		//		render.JSON(404, model.NewError(model.ERR_REQUEST_FAILED, err.Error()))
		//	} else {
		//		render.JSON(200, users)
		//	}
		//} else {
		//	var nusers []*model.User

		//	pname := query.Get("name")
		//	if !utils.IsEmpty(pname) {
		//		names := strings.Split(pname, ",")
		//		nusers, _ = model.GetUsersByName(db, names)
		//	}

		//	var iusers []*model.User
		//	pid := query.Get("id")
		//	if !utils.IsEmpty(pid) {
		//		var ids []int
		//		for _, v := range strings.Split(pid, ",") {
		//			if id, err := strconv.Atoi(v); err == nil {
		//				var exist bool = false
		//				for _, u := range nusers {
		//					if u.ID == id {
		//						exist = true
		//						break
		//					}
		//				}

		//				if !exist {
		//					ids = append(ids, id)
		//				}
		//			}
		//		}

		//		iusers, _ = model.GetUsersByID(db, ids)
		//	}

		//	users := make([]*model.User, len(iusers)+len(nusers))
		//	pos := copy(users, iusers)
		//	copy(users[pos:], nusers)

		//	render.JSON(200, users)
		//}
	}
}


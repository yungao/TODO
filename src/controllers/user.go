package controllers

import (
	"log"
	"net/http"
	"strconv"
	//"strings"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	//"github.com/martini-contrib/binding"

	model "models"
	utils "utils"
)

const (
	// request failed
	ERR_REQUEST_FAILED = 40100

	// ERR_INVALID_NAME    =   40011
	// ERR_INVALID_PWD     =   40012
	ERR_NAME_EXIST = 40013
)

/**
* Create a new user, to user register
 */
func CreateUser(session sessions.Session, user model.User, db *gorp.DbMap, render render.Render) {
	log.Println("Create User {Name: " + user.Name + ", Pwd: " + user.Pwd + "}")

	err := db.Insert(&user)
	if err == nil {
		render.JSON(201, &user)
	} else {
		log.Printf("Create user error: %s", err.Error())
		erp := model.Error{Code: ERR_NAME_EXIST, Msg: "Username already exists!"}
		render.JSON(422, erp)
	}
}

/**
* List user information
 */
func ListUser(session sessions.Session, db *gorp.DbMap, params martini.Params, render render.Render, request *http.Request) {
	query := request.URL.Query()
	name := query.Get("name")

	if !utils.IsEmpty(name) {
		//_, err := utils.ParseSession(session, render)
		//if err == nil {
		var user = model.User{}
		err := db.SelectOne(&user, "SELECT * FROM user WHERE name=?", name)
		if err != nil {
			render.JSON(200, "{}")
		} else {
			render.JSON(200, user)
		}
		//}
	} else {
		erp := model.Error{Code: ERR_REQUEST_FAILED, Msg: "Not Found!"}
		render.JSON(404, erp)
	}
}

/**
* user login
*
* Code:
*           11010: login failed
*           11011: user does not exist
*           11012: user password error
 */
func Login(session sessions.Session, user model.User, db *gorp.DbMap, render render.Render) {
	var dbUser = model.User{}
	err := db.SelectOne(&dbUser, "SELECT * FROM user WHERE name=?", user.Name)
	if err != nil {
		log.Printf("Login error: User[%s] does not exist", user.Name)
		erp := model.Error{Code: 11011, Msg: "User does not exist!"}
		render.JSON(422, erp)
	} else {
		if dbUser.Pwd != user.Pwd {
			log.Printf("Login error: User[%s]'s password[%s] error", user.Name, user.Pwd)
			erp := model.Error{Code: 11012, Msg: "User password error!"}
			render.JSON(422, erp)
		} else {
			session.Set("ID", dbUser.ID)
			log.Printf("Login Session: %s]", session)
			render.JSON(200, dbUser)
		}
	}
}

/**
* Get user information
*
* Code:    11011: user does not exist
 */
func GetUser(db *gorp.DbMap, params martini.Params, render render.Render) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err)
	}

	user, err := db.Get(model.User{}, id)
	if err == nil {
		render.JSON(200, user)
	} else {
		log.Printf("Get user error: %s", err.Error())
		erp := model.Error{Code: 11011, Msg: "User does not exist!"}
		render.JSON(404, erp)
	}
}

/**
* Delete user
*
* Code:    11021: delete user failed
 */
func DeleteUser(db *gorp.DbMap, params martini.Params, render render.Render) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err)
	}

	_, err = db.Delete(&model.User{ID: id})
	if err == nil {
		render.JSON(204, "No Content")
	} else {
		log.Printf("Delete user error: %s", err.Error())
		erp := model.Error{Code: 11021, Msg: "Delete user failed!"}
		render.JSON(404, erp)
	}
}

/**
* Change user password
*
* Code:    11041: changed password failed
 */
func changeUserPwd(db *gorp.DbMap, render render.Render, id int, pwd string) {
	_, err := db.Exec("UPDATE user SET pwd = ? WHERE id = ?", pwd, id)
	if err != nil {
		log.Printf("Change password error: %s", err.Error())
		erp := model.Error{Code: 11041, Msg: "Change user password failed!"}
		render.JSON(404, erp)
	} else {
		render.JSON(201, "OK")
	}
}

func UpdateUser(db *gorp.DbMap, params martini.Params, render render.Render, request *http.Request) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err)
	}

	pwd := request.FormValue("pwd")
	if pwd != "" {
		changeUserPwd(db, render, id, pwd)
		return
	}

	render.JSON(400, "Invaild Action")
}

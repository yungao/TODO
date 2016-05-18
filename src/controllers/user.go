package controllers

import (
	"log"
	"net/http"
	"strconv"
	// "regexp"
	// "strings"

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

	ERR_INVALID_DATA = 40011
	ERR_NAME_EXIST   = 40013
)

/**
* Create a new user, to user register
 */
func CreateUser(session sessions.Session, user model.User, db *gorp.DbMap, render render.Render) {
	log.Println("Create User: ", user.String())
	id, err := utils.ParseSession(session, render)
	if err == nil {
		_, err := model.DetermineAdmin(db, id)
		if err != nil {
			render.JSON(403, model.NewError(ERR_REQUEST_FAILED, err.Error()))
			return
		}

		// check name
		if err = model.VerifyName(user.Name); err != nil {
			render.JSON(422, model.NewError(ERR_INVALID_DATA, err.Error()))
			return
		}

		// check password
		if err = model.VerifyPassword(user.Pwd); err != nil {
			render.JSON(422, model.NewError(ERR_INVALID_DATA, err.Error()))
			return
		}

		// check email
		if user.Email != "" {
			if err = model.VerifyEmail(user.Email); err != nil {
				render.JSON(422, model.NewError(ERR_INVALID_DATA, err.Error()))
				return
			}
		}

		user.CreatorID = id
		err = db.Insert(&user)
		if err == nil {
			render.JSON(201, &user)
		} else {
			log.Printf("Create user error: %s", err.Error())
			erp := model.Error{Code: ERR_NAME_EXIST, Msg: "Name already exists!"}
			render.JSON(422, erp)
		}
	}
}

/*
* Update user info
 */
func UpdateUser(session sessions.Session, db *gorp.DbMap, params martini.Params, render render.Render, request *http.Request) {
	id, err := utils.ParseSession(session, render)
	if err == nil { // has login
		auth := request.FormValue("auth")
		active := request.FormValue("active")

		var needsAdmin bool = false
		if auth != "" || active != "" {
			needsAdmin = true
		}

		uid, err := strconv.Atoi(params["id"])
		if err != nil { // none ID param
			uid = id // update login user
		} else {
			if uid != id {
				needsAdmin = true
			}
		}

		if needsAdmin {
			_, err := model.DetermineAdmin(db, id)
			if err != nil {
				render.JSON(403, model.NewError(ERR_REQUEST_FAILED, err.Error()))
				return
			}
		}

		user, err := model.GetUserByID(db, uid)
		if err != nil { // user does not exist
			render.JSON(403, model.NewError(ERR_REQUEST_FAILED, "User does not exist!"))
			return
		}

		pwd := request.FormValue("pwd")
		if pwd != "" {
			// check password
			if err = model.VerifyPassword(pwd); err != nil {
				render.JSON(422, model.NewError(ERR_INVALID_DATA, err.Error()))
				return
			} else {
				user.Pwd = pwd
			}
		}

		nickname := request.FormValue("nickname")
		if nickname != "" {
			if err = model.VerifyNickName(nickname); err != nil {
				render.JSON(422, model.NewError(ERR_INVALID_DATA, err.Error()))
				return
			} else {
				user.Nickname = nickname
			}
		}

		email := request.FormValue("email")
		if email != "" {
			if err = model.VerifyEmail(email); err != nil {
				render.JSON(422, model.NewError(ERR_INVALID_DATA, err.Error()))
				return
			} else {
				user.Email = email
			}
		}

		if auth != "" {
			if i, err := strconv.ParseInt(auth, 10, 8); err == nil {
				user.Authority = int8(i)
			} else {
				render.JSON(422, model.NewError(ERR_INVALID_DATA, err.Error()))
				return
			}
		}

		if active != "" {
			if i, err := strconv.ParseInt(active, 10, 8); err == nil {
				user.Active = int8(i)
			} else {
				render.JSON(422, model.NewError(ERR_INVALID_DATA, err.Error()))
				return
			}
		}

		_, err = db.Update(user)
		if err == nil {
			render.JSON(201, user)
		} else {
			render.JSON(422, model.NewError(ERR_REQUEST_FAILED, err.Error()))
		}
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
		log.Printf("Delete tb_users error: %s", err.Error())
		erp := model.Error{Code: 11021, Msg: "Delete user failed!"}
		render.JSON(404, erp)
	}
}

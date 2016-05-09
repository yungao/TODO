package controller

import (
    "strconv"
    "log"
    "net/http"
    //"strings"

    "github.com/go-martini/martini"
    "github.com/coopernurse/gorp"
    "github.com/martini-contrib/render"
    //"github.com/martini-contrib/binding"

    "model"
)

/**
* create a new user
*
* errno:    11002: name is nil
*           11003: pwd is nil
*           11001: add failed
*/
func CreateUser(user model.User, db *gorp.DbMap, render render.Render) {
    log.Println("Add User {Name: " + user.Name + ", Pwd: " + user.Pwd + "}")
    if user.Name == "" {
        erp := model.Error{Errno: 11002, Msg: "User name can not be Null!"}
        // erp.Error();
        // render.JSON(422, map[string]interface{}{"result": false, "errno": 11001, "msg": "User name can not be Null!"})
        render.JSON(422, erp)
        return
    }
    if user.Pwd == "" {
        erp := model.Error{Errno: 11003, Msg: "User password can not be Null!"}
        render.JSON(422, erp)
        return
    }

    err := db.Insert(&user)
    if err == nil {
        render.JSON(201, &user)
    } else {
        log.Printf("Create user error: %s", err.Error())
        erp := model.Error{Errno: 11001, Msg: "Create user failed!"}
        render.JSON(422, erp)
    }
}

/**
* Get user information
*
* errno:    11011: user does not exist
*/
func GetUser(db *gorp.DbMap, params martini.Params, render render.Render) {
    id, err := strconv.ParseInt(params["id"], 0, 64)
    if err != nil {
        panic(err)
    }

    user, err := db.Get(model.User{}, id)
    if err == nil {
        render.JSON(200, user)
    } else {
        log.Printf("Get user error: %s", err.Error())
        erp := model.Error{Errno: 11011, Msg: "User does not exist!"}
        render.JSON(404, erp)
    }
}

/**
* Delete user
*
* errno:    11021: delete user failed
*/
func DeleteUser(db *gorp.DbMap, params martini.Params, render render.Render) {
    id, err := strconv.ParseInt(params["id"], 0, 64)
    if err != nil {
        panic(err)
    }

    _, err = db.Delete(&model.User{ID: id})
    if err == nil {
        render.JSON(204, "No Content")
    } else {
        log.Printf("Delete user error: %s", err.Error())
        erp := model.Error{Errno: 11021, Msg: "Delete user failed!"}
        render.JSON(404, erp)
    }
}

/**
* List all user information
*
* errno:    11031: list user failed
*/
func ListUser(db *gorp.DbMap, params martini.Params, render render.Render, request *http.Request) {
    query := request.URL.Query()
    var limit, offset string

    if query.Get("limit") != "" {
        limit = " LIMIT " + query.Get("limit")
    }
    if query.Get("offset") != "" {
        offset = " OFFSET " + query.Get("offset")
    }

    var users []model.User
    _, err := db.Select(&users, "SELECT * FROM user" + limit + offset)
    if err == nil {
        render.JSON(200, users)
    } else {
        log.Printf("List user error: %s", err.Error())
        erp := model.Error{Errno: 11031, Msg: "List user failed!"}
        render.JSON(404, erp)
    }
}

/**
* Change user password
*
* errno:    11041: changed password failed
*/
func changeUserPwd(db *gorp.DbMap, render render.Render, id int64, pwd string) {
    _, err := db.Exec("UPDATE user SET pwd = ? WHERE id = ?", pwd, id)
    if err != nil {
        log.Printf("Change password error: %s", err.Error())
        erp := model.Error{Errno: 11041, Msg: "Change user password failed!"}
        render.JSON(404, erp)
    } else {
        render.JSON(201, "OK")
    }
}

func UpdateUser(db *gorp.DbMap, params martini.Params, render render.Render, request *http.Request) {
    id, err := strconv.ParseInt(params["id"], 0, 64)
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


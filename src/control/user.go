package control

import (
    "strconv"
    "log"

    "github.com/go-martini/martini"
    "github.com/coopernurse/gorp"
    "github.com/martini-contrib/render"
    //"github.com/martini-contrib/binding"

    "model"
)

/**
* create a new user
*
* errno:    11001: name is nil
*           11002: pwd is nil
*           11003: add failed
*/
func CreateUser(user model.User, db *gorp.DbMap, render render.Render) {
    log.Println("Add User {Name: " + user.Name + ", Pwd: " + user.Pwd + "}")
    if user.Name == "" {
        rp := model.Response{Result: false, Errno: 11001, Msg: "User name can not be Null!", Data: nil}
        // rp.Response();
        // render.JSON(422, map[string]interface{}{"result": false, "errno": 11001, "msg": "User name can not be Null!"})
        render.JSON(422, &rp)
        return
    }
    if user.Pwd == "" {
        rp := model.Response{Result: false, Errno: 11002, Msg: "User password can not be Null!", Data: nil}
        render.JSON(422, &rp)
        return
    }

    err := db.Insert(&user)
    if err == nil {
        rp := model.Response{Result: true, Errno: 0, Msg: "", Data: user}
        render.JSON(201, &rp)
    } else {
        rp := model.Response{Result: false, Errno: 11003, Msg: "Add user failed!", Data: nil}
        render.JSON(422, &rp)
    }
}

func GetUser(db *gorp.DbMap, params martini.Params, render render.Render) {
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        panic(err)
    }
    log.Printf("Get User {ID: %d}", id)

    user, err := db.Get(model.User{}, id)
    if err == nil {
        rp := model.Response{Result: true, Errno: 0, Msg: "", Data: user}
        render.JSON(200, &rp)
    } else {
        rp := model.Response{Result: false, Errno: 11004, Msg: "User does not exist!", Data: nil}
        render.JSON(404, &rp)
    }
}

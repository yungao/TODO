package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	// "regexp"
    "errors"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	//"github.com/martini-contrib/binding"

	"config"
	model "models"
	utils "utils"
)

/**
* Check user id enabled
*/
func CheckUserEnable(db *gorp.DbMap, uid int) (*model.User, int, error) {
    user, err := model.GetUserByID(db, uid)
    if err != nil {
        return nil, model.ERR_USER_NOT_FOUND, err
    } else {
        if !user.IsEnable() {
            return nil, model.ERR_USER_NOT_FOUND, errors.New("User is disabled!")
        }
    }

    return user, 0, nil
}

/**
* Create a new user, to user register
 */
func CreateUser(session sessions.Session, user model.User, db *gorp.DbMap, render render.Render) {
	log.Println("Create User: ", user.String())
	id, name, err := utils.ParseSession(session, render)
	if err == nil {
		_, err := model.DetermineAdmin(db, id, name)
		if err != nil {
			render.JSON(403, model.NewError(model.ERR_REQUEST_FAILED, err.Error()))
			return
		}

		// check name
		if err = model.VerifyName(user.Name); err != nil {
			render.JSON(422, model.NewError(model.ERR_INVALID_DATA, err.Error()))
			return
		}

		// check password
		if err = model.VerifyPassword(user.Pwd); err != nil {
			render.JSON(422, model.NewError(model.ERR_INVALID_DATA, err.Error()))
			return
		}

		// check email
		if user.Email != "" {
			if err = model.VerifyEmail(user.Email); err != nil {
				render.JSON(422, model.NewError(model.ERR_INVALID_DATA, err.Error()))
				return
			}
		}

		user.CreatorID = id
		err = db.Insert(&user)
		if err == nil {
			render.JSON(201, &user)
		} else {
			log.Printf("Create user error: %s", err.Error())
			erp := model.Error{Code: model.ERR_NAME_EXIST, Msg: "Name already exists!"}
			render.JSON(422, erp)
		}
	}
}

/*
* Update user info
 */
func UpdateUser(session sessions.Session, db *gorp.DbMap, params martini.Params, render render.Render, request *http.Request) {
	log.Println("Requst request:", request)
	id, name, err := utils.ParseSession(session, render)
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
			_, err := model.DetermineAdmin(db, id, name)
			if err != nil {
				render.JSON(403, model.NewError(model.ERR_REQUEST_FAILED, err.Error()))
				return
			}
		}

		user, err := model.GetUserByID(db, uid)
		if err != nil { // user does not exist
			render.JSON(403, model.NewError(model.ERR_USER_NOT_FOUND, "User does not exist!"))
			return
		}

		pwd := request.FormValue("pwd")
		if pwd != "" {
			// check password
			if err = model.VerifyPassword(pwd); err != nil {
				render.JSON(422, model.NewError(model.ERR_INVALID_DATA, err.Error()))
				return
			} else {
				user.Pwd = utils.Base64Encode(pwd)
			}
		}

		nickname := request.FormValue("nickname")
		if nickname != "" {
			if err = model.VerifyNickName(nickname); err != nil {
				render.JSON(422, model.NewError(model.ERR_INVALID_DATA, err.Error()))
				return
			} else {
				user.Nickname = nickname
			}
		}

		email := request.FormValue("email")
		if email != "" {
			if err = model.VerifyEmail(email); err != nil {
				render.JSON(422, model.NewError(model.ERR_INVALID_DATA, err.Error()))
				return
			} else {
				user.Email = email
			}
		}

		if auth != "" {
			if i, err := strconv.ParseInt(auth, 10, 8); err == nil {
				user.Authority = int8(i)
			} else {
				render.JSON(422, model.NewError(model.ERR_INVALID_DATA, err.Error()))
				return
			}
		}

		if active != "" {
			if i, err := strconv.ParseInt(active, 10, 8); err == nil {
				user.Active = int8(i)
			} else {
				render.JSON(422, model.NewError(model.ERR_INVALID_DATA, err.Error()))
				return
			}
		}

		_, err = db.Update(user)
		if err == nil {
			render.JSON(201, user)
		} else {
			render.JSON(422, model.NewError(model.ERR_REQUEST_FAILED, err.Error()))
		}
	}
}

/**
* User login
 */
func Login(session sessions.Session, user model.User, db *gorp.DbMap, render render.Render) {
	var dbUser = model.User{}
	err := db.SelectOne(&dbUser, "SELECT * FROM "+config.TABLE_NAME_USERS+" WHERE name=?", user.Name)
	if err != nil {
		log.Printf("Login error: User[%s] does not exist", user.Name)
		render.JSON(422, model.NewError(model.ERR_USER_NOT_FOUND, fmt.Sprintf("User[%s] does not exist!", user.Name)))
	} else {
		if !dbUser.IsEnable() {
			render.JSON(403, model.NewError(model.ERR_USER_DISABLED, "The user is disabled!"))
		} else {
			pwd := utils.Base64Encode(user.Pwd)
			if dbUser.Pwd != pwd {
				log.Printf("Login error: User[%s]'s password[%s] error", user.Name, user.Pwd)
				erp := model.NewError(model.ERR_INVALID_DATA, fmt.Sprintf("User[%s]'s password[%s] error!", user.Name, user.Pwd))
				render.JSON(422, erp)
			} else {
				if dbUser.Active == 0 {
					dbUser.Active = 1
					db.Update(&dbUser)
				}

				s := fmt.Sprintf("%d:%s", dbUser.ID, dbUser.Name)
				session.Set("ID", s)
				log.Println("Login Session: ", session)
				render.JSON(201, dbUser)
			}
		}
	}
}

/**
* User logout
 */
func Logout(session sessions.Session, db *gorp.DbMap, params martini.Params, render render.Render) {
	_, _, err := utils.ParseSession(session, render)
	if err == nil { // has login
		session.Delete("ID")
		render.JSON(200, "Logout!")
	}
}

/**
* Get user information
 */
func GetUser(session sessions.Session, db *gorp.DbMap, params martini.Params, render render.Render) {
	_, _, err := utils.ParseSession(session, render)
	if err == nil { // has login
		param := params["id"]
		id, err := strconv.Atoi(param)
		if err == nil { // get by ID
			user, err := model.GetUserByID(db, id)
			if err == nil {
				render.JSON(200, user)
				return
			}
		} else { // get by Name
			user, err := model.GetUserByName(db, param)
			if err == nil {
				render.JSON(200, user)
				return
			}
		}

		render.JSON(404, model.NewError(model.ERR_USER_NOT_FOUND, fmt.Sprintf("User[%s] does not exist!", param)))
	}
}

/**
* List users info
 */
func ListUsers(session sessions.Session, db *gorp.DbMap, params martini.Params, render render.Render, request *http.Request) {
	_, _, err := utils.ParseSession(session, render)
	if err == nil { // has login
		query := request.URL.Query()
		if len(query) == 0 { // list all users
			var users, err = model.GetAllUsers(db)
			if err != nil {
				render.JSON(404, model.NewError(model.ERR_REQUEST_FAILED, err.Error()))
			} else {
				render.JSON(200, users)
			}
		} else {
			var nusers []*model.User

			pname := query.Get("name")
			if !utils.IsEmpty(pname) {
				names := strings.Split(pname, ",")
				nusers, _ = model.GetUsersByName(db, names)
			}

			var iusers []*model.User
			pid := query.Get("id")
			if !utils.IsEmpty(pid) {
				var ids []int
				for _, v := range strings.Split(pid, ",") {
					if id, err := strconv.Atoi(v); err == nil {
						var exist bool = false
						for _, u := range nusers {
							if u.ID == id {
								exist = true
								break
							}
						}

						if !exist {
							ids = append(ids, id)
						}
					}
				}

				iusers, _ = model.GetUsersByID(db, ids)
			}

			users := make([]*model.User, len(iusers)+len(nusers))
			pos := copy(users, iusers)
			copy(users[pos:], nusers)

			render.JSON(200, users)
		}
	}
}

/**
* Delete user
 */
func DeleteUser(session sessions.Session, db *gorp.DbMap, params martini.Params, render render.Render) {
	id, name, err := utils.ParseSession(session, render)
	if err == nil { // has login
		_, err := model.DetermineAdmin(db, id, name)
		if err == nil { // login user is admin
			uid, err := strconv.Atoi(params["id"])
			if err == nil {
				if id == uid {
					render.JSON(403, model.NewError(model.ERR_REQUEST_FAILED, "Can not delete admin!"))
				} else {
					count, err := db.Delete(&model.User{ID: uid})
					if err == nil {
						if count == 1 {
							render.JSON(204, "Delete success!")
						} else {
							log.Printf("Delete user[%d] does not exist: %s", uid)
							render.JSON(404, model.NewError(model.ERR_USER_NOT_FOUND, fmt.Sprintf("User[%d] does not exist!", uid)))
						}
					} else {
						log.Printf("Delete tb_users error: %s", err.Error())
						render.JSON(404, model.NewError(model.ERR_REQUEST_FAILED, err.Error()))
					}
				}
			} else {
				render.JSON(404, model.NewError(model.ERR_INVALID_DATA, "Invalid Data!"))
			}
		} else {
			render.JSON(403, model.NewError(model.ERR_REQUEST_FAILED, err.Error()))
		}
	}
}

func IsLogin(session sessions.Session, db *gorp.DbMap, params martini.Params, render render.Render) {
	id, _, err := utils.ParseSession(session, render)
	if err == nil { // has login
		user, err := model.GetUserByID(db, id)
		if err == nil {
			render.JSON(200, user)
			return
		}
	}
}

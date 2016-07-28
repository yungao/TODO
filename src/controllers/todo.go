package controllers

import (
	"fmt"
    "time"
	"log"
	"net/http"
	 "strconv"
	// "strings"
	// "regexp"
    "errors"

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
* Check todo enabled
*/
func CheckTodoEnable(db *gorp.DbMap, tid int) (*model.Todo, int, error) {
    todo, err := model.GetTodoByID(db, tid)
    if err != nil {
        return nil, model.ERR_TODO_NOT_FOUND, err
    } else {
        if !todo.IsEnable() {
            return nil, model.ERR_TODO_DISABLED, errors.New("Todo is disabled!")
        }
    }

    return todo, 0, nil
}

/**
* Create a new Todo
 */
func CreateTodo(session sessions.Session, todo model.Todo, db *gorp.DbMap, render render.Render, request *http.Request) {
	log.Println("Create Todo: ", todo.String())
	id, _, err := utils.ParseSession(session, render)
    if err == nil {
        user, err := model.GetUserByID(db, id)
        if err == nil {
            if user.IsEnable() {
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
                    agent := utils.ParseUserAgent(request)
                    err = CreateTodoProcess(db, &todo, user, agent)
                    if err == nil {
                        render.JSON(201, &todo)
                    } else {
                        log.Printf("Create todo process error: %s", err.Error())
                        erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: err.Error()}
                        render.JSON(422, erp)
                    }
                } else {
                    log.Printf("Create todo error: %s", err.Error())
                    erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: err.Error()}
                    render.JSON(422, erp)
                }
            } else {
                erp := model.Error{Code: model.ERR_USER_DISABLED, Msg: "User is disabled!"}
                render.JSON(403, erp)
            }
        } else {
            erp := model.Error{Code: model.ERR_USER_NOT_FOUND, Msg: err.Error()}
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
                todo.Creator, _ = model.GetUserByID(db, todo.CreatorID)
                todo.Partners, _ = model.GetPartnersByTodoID(db, todo.ID)
                // query todo partner user info
                for _, partner := range todo.Partners {
                    partner.Partner, _ = model.GetUserByID(db, partner.PartnerID)
                }
                // query todo partner creator user info
                for _, partner := range todo.Partners {
                    exist := false
                    for _, p := range todo.Partners {
                        if partner.UserID == p.PartnerID {
                            partner.Creator = p.Partner
                            exist = true
                            break;
                        }
                    }

                    if !exist {
                        partner.Creator, _ = model.GetUserByID(db, partner.UserID)
                    }
                }

                processes, err := model.GetTodoProcesses(db, todo.ID)
                todo.Processes = processes
                for _, process := range todo.Processes {
                    exist := false
                    for _, partner := range todo.Partners {
                        if process.CreatorID == partner.UserID {
                            process.Creator = partner.Creator
                            exist = true
                            break;
                        }
                    }

                    if !exist {
                        process.Creator, _ = model.GetUserByID(db, process.CreatorID)
                    }
                }
                if err == nil {
                    render.JSON(200, todo)
                    return
                }
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
            for _, todo := range todos {
                todo.Creator, _ = model.GetUserByID(db, todo.CreatorID);
            }

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

func doUpdateTodo(db *gorp.DbMap, user *model.User, todo *model.Todo, render render.Render, request *http.Request) {
    var hasPermission = (todo.Type == 0)
    if !hasPermission {
        hasPermission = (user.ID == todo.CreatorID)
    }

    partners, _ := model.GetPartnersByTodoID(db, todo.ID)
    for _, p := range partners {
        if p.PartnerID == user.ID {
            hasPermission = true
            break;
        }
    }

    if hasPermission {
        content := request.FormValue("content")
        if len(content) > 0 {
            agent := utils.ParseUserAgent(request)
            err := UpdateTodoProcess(db, todo, user, content, agent)
            if err == nil {
                render.JSON(201, &todo)
            } else {
                log.Printf("Update todo error: %s", err.Error())
                erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: err.Error()}
                render.JSON(422, erp)
            }
        } else {
            log.Print("Update todo error: No content to update")
            erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: "No content to update!"}
            render.JSON(422, erp)
        }
    } else {
        log.Print("Update todo error: No permission")
        erp := model.Error{Code: model.ERR_TODO_USER_PERM, Msg: "Do not have permission to update!"}
        render.JSON(422, erp)
    }
}

func doModifyPriority(db *gorp.DbMap, user *model.User, todo *model.Todo, render render.Render, request *http.Request) {
    if user.ID == todo.CreatorID {
        p, err := strconv.ParseInt(request.FormValue("priority"), 10, 8)
        if err == nil && todo.Priority != int8(p) {
            priority := int8(p)
            todo.Priority = priority
            _, err = db.Update(todo)
            if err == nil {
                agent := utils.ParseUserAgent(request)
                err := ModifyPriorityProcess(db, todo, user, priority, agent)
                if err == nil {
                    render.JSON(201, &todo)
                } else {
                    log.Printf("Update todo error: %s", err.Error())
                    erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: err.Error()}
                    render.JSON(422, erp)
                }
            } else {
                log.Printf("Update todo error: %s", err.Error())
                erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: err.Error()}
                render.JSON(422, erp)
            }
        } else {
            log.Print("Update todo error: No priority to update")
            erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: "No priority to update!"}
            render.JSON(422, erp)
        }
    } else {
        log.Print("Update todo error: No permission")
        erp := model.Error{Code: model.ERR_TODO_USER_PERM, Msg: "Do not have permission to update!"}
        render.JSON(422, erp)
    }
}

func doModifyLimit(db *gorp.DbMap, user *model.User, todo *model.Todo, render render.Render, request *http.Request) {
    if user.ID == todo.CreatorID {
        p, err := strconv.ParseInt(request.FormValue("limit"), 10, 64)
        if err == nil && todo.LimitAt != int64(p) {
            limit := int64(p)
            todo.LimitAt = limit
            _, err = db.Update(todo)
            if err == nil {
                agent := utils.ParseUserAgent(request)
                err := ModifyLimitProcess(db, todo, user, limit, agent)
                if err == nil {
                    render.JSON(201, &todo)
                } else {
                    log.Printf("Update todo error: %s", err.Error())
                    erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: err.Error()}
                    render.JSON(422, erp)
                }
            } else {
                log.Printf("Update todo error: %s", err.Error())
                erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: err.Error()}
                render.JSON(422, erp)
            }
        } else {
            log.Print("Update todo error: No limit to update")
            erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: "No limit to update!"}
            render.JSON(422, erp)
        }
    } else {
        log.Print("Update todo error: No permission")
        erp := model.Error{Code: model.ERR_TODO_USER_PERM, Msg: "Do not have permission to update!"}
        render.JSON(422, erp)
    }
}

/*
* Update todo info
 */
 func UpdateTodo(session sessions.Session, db *gorp.DbMap, params martini.Params, render render.Render, request *http.Request) {
     log.Println("Requst request:", request)
     uid, _, err := utils.ParseSession(session, render)
     if err == nil { // has login
         user, code, err := CheckUserEnable(db, uid)
         if err == nil {
             tid, err := strconv.Atoi(params["id"])
             if err == nil {
                 todo, code, err := CheckTodoEnable(db, tid)
                 todo.ProcCount += 1
                 if err == nil {
                     action, err := strconv.Atoi(request.FormValue("action"))
                     if err == nil {
                         switch action {
                             case 0: // 完成任务

                             case 2: // 添加参与者

                             case 3: // 指派执行者

                             case 4: // 添加TAG

                             case 5: // 修改结束时间
                                doModifyLimit(db, user, todo, render, request)
                             case 6: // 修改紧急度
                                doModifyPriority(db, user, todo, render, request)
                             case 7: // 更新任务
                                doUpdateTodo(db, user, todo, render, request)
                             case -1: // 撤销完成
                         }
                     } else {
                         log.Print("Update todo error: No action")
                         erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: "No param of Update Action!"}
                         render.JSON(422, erp)
                     }
                 } else {
                     log.Printf("Update todo error: %s", err.Error())
                     erp := model.Error{Code: code, Msg: err.Error()}
                     render.JSON(422, erp)
                 }
             } else {
                 log.Print("Update todo error: No todo id")
                 erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: "No param of Todo ID!"}
                 render.JSON(422, erp)
             }
         } else {
             erp := model.Error{Code: code, Msg: err.Error()}
             render.JSON(403, erp)
         }
     }
 }

package controllers

import (
	"log"
	"net/http"

	"github.com/coopernurse/gorp"
	// "github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	//"github.com/martini-contrib/binding"

	// "config"
	model "models"
	utils "utils"
)

/**
* Create a new Tag
 */
 func CreateTag(session sessions.Session, tag model.Tag, db *gorp.DbMap, render render.Render) {
     log.Println("Create Tag: ", tag.String())
     uid, _, err := utils.ParseSession(session, render)
     if err == nil {
         _, code, err := CheckUserEnable(db, uid)
         if err == nil {
             tag.CreatorID = uid
             err = db.Insert(&tag)
             if err == nil {
                 render.JSON(201, &tag)
             } else {
                 log.Printf("Create tag error: %s", err.Error())
                 erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: err.Error()}
                 render.JSON(422, erp)
             }
         } else {
             erp := model.Error{Code: code, Msg: err.Error()}
             render.JSON(403, erp)
         }
     }
 }

/**
* Add Todo Tag
 */
 func AddTodoTag(session sessions.Session, ttag model.TodoTag, db *gorp.DbMap, render render.Render, request *http.Request) {
     log.Println("Add Todo Tag: ", ttag.String())
     uid, _, err := utils.ParseSession(session, render)
     if err == nil {
         user, code, err := CheckUserEnable(db, uid)
         if err == nil {
             todo, code, err := CheckTodoEnable(db, ttag.TodoID)
             if err == nil {
                 isPartner, err := IsTodoPartner(db, todo, user)
                 if err == nil && isPartner {
                     ttag.CreatorID = uid
                     err = db.Insert(&ttag)
                     if err == nil {
                         render.JSON(201, &ttag)
                     } else {
                         log.Printf("Add todo tag error: %s", err.Error())
                         erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: err.Error()}
                         render.JSON(422, erp)
                     }
                 } else {
                     log.Print("Add todo tag error: has no permission")
                     erp := model.Error{Code: model.ERR_TODO_USER_PERM, Msg: "No permission!"}
                     render.JSON(422, erp)
                 }
             } else {
                 log.Printf("Add todo tag error: %s", err.Error())
                 erp := model.Error{Code: code, Msg: err.Error()}
                 render.JSON(422, erp)
             }
         } else {
             erp := model.Error{Code: code, Msg: err.Error()}
             render.JSON(403, erp)
         }
     }
 }

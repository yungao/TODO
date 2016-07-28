package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coopernurse/gorp"
	// "github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	//"github.com/martini-contrib/binding"

	"config"
	model "models"
	utils "utils"
)

// func isPartnerExist(db *gorp.DbMap, tid int, pid int) (bool, error) {
// 	rows, err := db.Select(model.Partner{}, fmt.Sprintf("SELECT * FROM %s WHERE %s.todoid=%d AND %s.pid=%d", config.TABLE_NAME_TODO_PARTNER, config.TABLE_NAME_TODO_PARTNER, tid, config.TABLE_NAME_TODO_PARTNER, pid))
//     return len(rows) > 0, err
// }

/**
* Add Todo Partner
 */
func AddPartner(session sessions.Session, partner model.Partner, db *gorp.DbMap, render render.Render, request *http.Request) {
	log.Println("Add partner: ", partner.String())
	uid, _, err := utils.ParseSession(session, render)
    if err == nil {
        user, code, err := CheckUserEnable(db, uid)
        if err == nil {
            todo, code, err := CheckTodoEnable(db, partner.TodoID)
            if err == nil {
                if uid != partner.PartnerID {
                    partner.UserID = uid
                    err = db.Insert(&partner)
                    err = AddPartnerProcess(db, todo, user, partner.PartnerID, utils.ParseUserAgent(request))
                    if err == nil {
                        render.JSON(201, &partner)
                    } else {
                        log.Printf("Add todo partner error: %s", err.Error())
                        erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: err.Error()}
                        render.JSON(422, erp)
                    }
                } else {
                    log.Printf("Add todo partner error: %s", err.Error())
                    erp := model.Error{Code: model.ERR_REQUEST_FAILED, Msg: err.Error()}
                    render.JSON(422, erp)
                }
            } else {
                log.Printf("Add todo partner error: %s", err.Error())
                erp := model.Error{Code: code, Msg: err.Error()}
                render.JSON(422, erp)
            }
        } else {
            erp := model.Error{Code: code, Msg: err.Error()}
            render.JSON(403, erp)
        }
    }
}

func GetPartnersByTodoID(db *gorp.DbMap, tid int) ([]*model.Partner, error) {
    var partners []*model.Partner
	_, err := db.Select(partners, fmt.Sprintf("SELECT * FROM %s WHERE %s.todoid=%d", config.TABLE_NAME_TODO_PARTNER, config.TABLE_NAME_TODO_PARTNER, tid))

    return partners, err
}

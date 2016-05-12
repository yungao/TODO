package models

import (
    "log"
    "time"
    "github.com/coopernurse/gorp"
    "github.com/martini-contrib/sessionauth"

    "config"
)

type User struct {
    ID          int64       `json:"id"          db:"id"`
    Name        string      `json:"name"        db:"name"       form:"name"     binding:"required"`
    Pwd         string      `json:"-"           db:"pwd"        form:"pwd"      binding:"required"`
    CreateAt    int64       `json:"create"      db:"create"`
    UpdateAt    int64       `json:"update"      db:"update"`
    authed      bool        `json:"-"           db:"-"          form:"-"`
}

type Users struct {
    Collection []User `json:"users"`
}

func (user *User) PreInsert(s gorp.SqlExecutor) error {
    user.CreateAt = time.Now().Unix()
    user.UpdateAt = user.CreateAt
    return nil
}

func (user *User) PreUpdate(s gorp.SqlExecutor) error {
    user.UpdateAt = time.Now().Unix()
    return nil
}

//********************** For sessionauth *****************************
func GenerateAnonymousUser() sessionauth.User {
    return &User{}
}

func (uer *User) Login() {
    uer.authed = true
}

// Logout will preform any actions that are required to completely logout a user.
func (u *User) Logout() {
    u.authed = false
}

func (user *User) IsAuthenticated() bool {
    return user.authed
}

func (user *User) UniqueId() interface{} {
    log.Printf("UniqueId: %s", user)
    return user.ID
}

// GetById will populate a user object from a database model with a matching id.
func (user *User) GetById(id interface{}) error {
    db := config.GetDBMap()
    err := db.SelectOne(user, "SELECT * FROM user WHERE id = $1", id)
    if err != nil {
        log.Printf("GetById error: %s", err.Error())
        return err
    }

    return nil
}
//**********************************************************************

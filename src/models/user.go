package models

import (
    //"log"
    "time"
    "github.com/coopernurse/gorp"
)

type User struct {
    ID          int64       `json:"id"          db:"id"`
    Name        string      `json:"name"        db:"name"       form:"name"     binding:"required"`
    Pwd         string      `json:"-"           db:"pwd"        form:"pwd"      binding:"required"`
    CreateAt    int64       `json:"create"      db:"create"`
    UpdateAt    int64       `json:"update"      db:"update"`
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

package models

type UserSession struct {
    Id            int64  `form:"id" db:"id"`
    Username      string `form:"name" db:"username"`
    Password      string `form:"password" db:"password"`
    authenticated bool   `form:"-" db:"-"`
}

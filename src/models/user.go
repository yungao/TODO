package models

import (
	"config"
	"github.com/coopernurse/gorp"
	"log"
	"time"
)

type User struct {
	ID       int    `db:"id"        json:"id"`
	Name     string `db:"name"      json:"name"         form:"name"     binding:"required"`
	Pwd      string `db:"pwd"       json:"-"            form:"pwd"      binding:"required"`
	Nickname string `db:"nickname"  json:"nickname"     form:"nickname"`
	Email    string `db:"email"     json:"email"        form:"email"`
	/* Authority:
	 *      -1:  super admin (create default)
	 *      0:   normal
	 *      1:   admin
	 */
	Authority int8  `db:"auth"      json:"auth"`
	CreatorID int   `db:"uid"       json:"-"`
	CreateAt  int64 `db:"create"    json:"create"`
	UpdateAt  int64 `db:"update"    json:"update"`
	/* Active:
	 *       -2: disabled
	 *       -1: deleted
	 *       0:  unused
	 *       1:  normal
	 */
	Active int8 `db:"active"    json:"active"`
}

type Users struct {
	Collection []User `json:"users"`
}

// Create user table if not exist
func CreateUserTable(db *gorp.DbMap) {
	tb := db.AddTableWithName(User{}, config.TABLE_NAME_USER)
	tb.SetKeys(true, "id")
	tb.ColMap("name").SetMaxSize(20).SetUnique(true).SetNotNull(true)
	tb.ColMap("pwd").SetMaxSize(20).SetNotNull(true)
	tb.ColMap("nickname").SetMaxSize(100)
	tb.ColMap("email").SetMaxSize(100)
	tb.ColMap("auth").SetNotNull(true)
	tb.ColMap("uid").SetNotNull(true)
	tb.ColMap("create").SetNotNull(true)
	tb.ColMap("update").SetNotNull(true)
	tb.ColMap("active").SetNotNull(true)

	err := db.CreateTablesIfNotExists()
	// db.DropTables()
	if err != nil {
		panic(err)
	}

	log.Printf(">>> Table[%s] created", config.TABLE_NAME_USER)
}

func (user *User) PreInsert(s gorp.SqlExecutor) error {
	user.CreateAt = time.Now().Unix()
	user.UpdateAt = user.CreateAt
	user.Active = 0
	return nil
}

func (user *User) PreUpdate(s gorp.SqlExecutor) error {
	user.UpdateAt = time.Now().Unix()
	return nil
}

type Partner struct {
	ID     int `db:"id"        json:"id"`
	TodoID int `db:"todoid"    json:"todoid"   form:"todoid"   binding:"required"`
	UserID int `db:"uid"       json:"uid"      form:"name"     binding:"required"`
	/* Active:
	 *       -1: deleted
	 *       1:  normal
	 */
	Active int8 `db:"active"    json:"active"`
}

type Partners struct {
	Collection []Partner `json:"partners"`
}

// Create partner table if not exist
func CreatePartnerTable(db *gorp.DbMap) {
	tb := db.AddTableWithName(Partner{}, config.TABLE_NAME_TODO_PARTNER)
	tb.SetKeys(true, "id")
	tb.ColMap("todoid").SetNotNull(true)
	tb.ColMap("uid").SetNotNull(true)
	tb.SetUniqueTogether("todoid", "uid")
	tb.ColMap("active").SetNotNull(true)

	err := db.CreateTablesIfNotExists()
	// db.DropTables()
	if err != nil {
		panic(err)
	}

	log.Printf(">>> Table[%s] created", config.TABLE_NAME_TODO_PARTNER)
}

func (partner *Partner) PreInsert(s gorp.SqlExecutor) error {
	partner.Active = 1
	return nil
}

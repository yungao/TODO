package models

import (
	"github.com/coopernurse/gorp"

	"fmt"
	"log"
	"time"

	"config"
	"utils"
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

func newAdmin() *User {
	return &User{
		Name:      "admin",
		Pwd:       "123456",
		Nickname:  "管理员",
		Email:     "admin@todo.com",
		Authority: -1,
		CreatorID: 0,
		CreateAt:  time.Now().Unix(),
		UpdateAt:  time.Now().Unix(),
		Active:    0,
	}
}

func (user *User) String() string {
	return fmt.Sprintf("{ID:%d, Name:%s, Pwd:%s, Nickname:%s, Email:%s, Authority:%d, CreatorID:%d, CreateAt:%d, UpdateAt:%d, Active:%d}", user.ID, user.Name, user.Pwd, user.Nickname, user.Email, user.Authority, user.CreatorID, user.CreateAt, user.UpdateAt, user.Active)
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

	// create super admin
	db.Insert(newAdmin())
	log.Println(">>> Admin user created")
}

func (user *User) PreInsert(s gorp.SqlExecutor) error {
	user.CreateAt = time.Now().Unix()
	user.UpdateAt = user.CreateAt
	user.Active = 0
	user.Pwd = utils.Base64Encode(user.Pwd)
	return nil
}

func (user *User) PreUpdate(s gorp.SqlExecutor) error {
	user.UpdateAt = time.Now().Unix()
	return nil
}

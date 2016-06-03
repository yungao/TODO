package models

import (
	"github.com/coopernurse/gorp"

	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"config"
	"utils"
)

const VALUE_ADMIN_AUTHORITY = -1

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

func (user *User) IsEnable() bool {
    return user.Active != -1
}

func (user *User) String() string {
	return fmt.Sprintf("{ID:%d, Name:%s, Pwd:%s, Nickname:%s, Email:%s, Authority:%d, CreatorID:%d, CreateAt:%d, UpdateAt:%d, Active:%d}", user.ID, user.Name, user.Pwd, user.Nickname, user.Email, user.Authority, user.CreatorID, user.CreateAt, user.UpdateAt, user.Active)
}

// Create user table if not exist
func CreateUserTable(db *gorp.DbMap) {
	tb := db.AddTableWithName(User{}, config.TABLE_NAME_USERS)
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
	log.Printf(">>> Table[%s] created", config.TABLE_NAME_USERS)

	// create super admin
	db.Insert(newAdmin())
	log.Println(">>> Admin user created")
}

// verify name
func VerifyName(name string) error {
	nlen := len(name)
	if nlen < 3 || nlen > 20 {
		return errors.New("Name must be 3-20 characters!")
	}

	if m, _ := regexp.MatchString("^[0-9a-zA-Z_]+$", name); !m {
		return errors.New("Name must be [0-9a-zA-Z_]!")
	}

	return nil
}

// verify password
func VerifyPassword(pwd string) error {
	plen := len(pwd)
	if plen < 3 || plen > 20 {
		return errors.New("Password must be 3-20 characters!")
	}

	if m, _ := regexp.MatchString("^[0-9a-zA-Z_]+$", pwd); !m {
		return errors.New("Password must be [0-9a-zA-Z_]!")
	}

	return nil
}

// verify email
func VerifyEmail(email string) error {
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, email); !m {
		return errors.New("Invaild email address!")
	}

	return nil
}

// verify nickname
func VerifyNickName(nick string) error {
	if len(nick) > 100 {
		return errors.New("Nickname must be less than 100 characters!")
	}

	return nil
}

/*
* To determine whether the user is an administrator
 */
func DetermineAdmin(db *gorp.DbMap, id int, name string) (bool, error) {
	if name == "admin" {
		return true, nil
	}

	// ret, err := db.Get(User{}, id)
	// if err != nil {
	// 	log.Printf("Login user does not exist: %s", err.Error())
	// 	return nil, errors.New("Login user does not exist!")
	// }

	// if user, ok := ret.(*User); ok && user.Authority == VALUE_ADMIN_AUTHORITY {
	// 	return user, nil
	// }
	return false, errors.New("Must login with admin!")
}

/*
* Query user info from database by ID
 */
func GetUserByID(db *gorp.DbMap, id int) (*User, error) {
	ret, err := db.Get(User{}, id)
	if err != nil {
		return nil, err
	}

	if user, ok := ret.(*User); ok {
		return user, nil
	}

	return nil, errors.New("User does not exist!")
}

/*
* Query multiple users info from database by ID
 */
func GetUsersByID(db *gorp.DbMap, ids []int) ([]*User, error) {
	var users []*User
	for _, id := range ids {
		ret, err := db.Get(User{}, id)
		if err == nil {
			if user, ok := ret.(*User); ok {
				users = append(users, user)
			}
		}
	}

	return users, nil
}

/*
* Query user info from database by Name
 */
func GetUserByName(db *gorp.DbMap, name string) (*User, error) {
	var user = User{}
	err := db.SelectOne(&user, "SELECT * FROM " + config.TABLE_NAME_USERS + " WHERE name = ?", name)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

/*
* Query users info from database by Name
 */
func GetUsersByName(db *gorp.DbMap, names []string) ([]*User, error) {
	var users []*User
	for _, name := range names {
		if user, err := GetUserByName(db, name); err == nil {
			users = append(users, user)
		}
	}

	return users, nil
}

/*
* Query all users info from database
 */
func GetAllUsers(db *gorp.DbMap) ([]*User, error) {
	var users []*User
	_, err := db.Select(&users, "SELECT * FROM "+config.TABLE_NAME_USERS)
	if err != nil {
		return nil, err
	}

	return users, nil
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

package models

import (
	"config"
	"github.com/coopernurse/gorp"
	"log"
	"time"
	"fmt"
    "errors"
)

type Todo struct {
	ID        int    `db:"id"        json:"id"`
	CreatorID int    `db:"uid"       json:"-"`
	Creator   *User  `db:"-"         json:"creator"`
	Name      string `db:"name"      json:"name"        form:"name"     binding:"required"`
	// Name      string `db:"name"      json:"name"        form:"name"`
	Content   string `db:"content"   json:"content"     form:"content"`
	Priority  int8   `db:"priority"  json:"priority"    form:"priority"`
	/* Type:
	 *       0:  public
	 *       1:  private
	 */
	Type     int8  `db:"type"      json:"type"        form:"type"`
	CreateAt int64 `db:"create"    json:"create"`
	LimitAt  int64 `db:"limit"     json:"limit"       form:"limit"    binding:"required"`
	UpdateAt int64 `db:"update"    json:"update"`
	/* Status:
	 *       0: undone
	 *       1: done
	 */
	Status    int8 `db:"status"    json:"status"`
	ProcCount int  `db:"pcount"    json:"pcount"`
	/* Active:
	 *       -1: deleted
	 *       0:  diabled
	 *       1:  normal
	 */
	Active     int8         `db:"active"    json:"active"`

	Processes  []*Process   `db:"-"         json:"processes"`
	Partners   []*Partner   `db:"-"         json:"partners"`
}

type Todos struct {
	Collection []Todo `json:"todos"`
}

// Create todo table if not exist
func CreateTodoTable(db *gorp.DbMap) {
	tb := db.AddTableWithName(Todo{}, config.TABLE_NAME_TODOS)
	tb.SetKeys(true, "id")
	tb.ColMap("uid").SetNotNull(true)
	tb.ColMap("name").SetMaxSize(100).SetNotNull(true)
	tb.ColMap("content").SetMaxSize(2048)
	tb.ColMap("priority").SetNotNull(true)
	tb.ColMap("type").SetNotNull(true)
	tb.ColMap("create").SetNotNull(true)
	tb.ColMap("limit").SetNotNull(true)
	tb.ColMap("update").SetNotNull(true)
	tb.ColMap("status").SetNotNull(true)
	tb.ColMap("pcount").SetNotNull(true)
	tb.ColMap("active").SetNotNull(true)

	err := db.CreateTablesIfNotExists()
	// db.DropTables()
	if err != nil {
		panic(err)
	}

	log.Printf(">>> Table[%s] created", config.TABLE_NAME_TODOS)
}

func (todo *Todo) String() string {
    return fmt.Sprintf("{ID:%d, CreatorID:%d, Name:%s, Content:%s, Priority:%d, Type:%d, CreateAt:%d, LimitAt:%d, UpdateAt:%d, Status:%d, ProcessCount:%d, Active:%d}", todo.ID, todo.CreatorID, todo.Name, todo.Content, todo.Priority, todo.Type, todo.CreateAt, todo.LimitAt, todo.UpdateAt, todo.Status, todo.ProcCount, todo.Active)
}

/*
* Query todo info from database by ID
 */
func GetTodoByID(db *gorp.DbMap, id int) (*Todo, error) {
	ret, err := db.Get(Todo{}, id)
	if err != nil {
		return nil, err
	}

	if todo, ok := ret.(*Todo); ok {
		return todo, nil
	}

	return nil, errors.New("Todo does not exist!")
}

/*
* Query all todos info from database
 */
func GetAllTodos(db *gorp.DbMap) ([]*Todo, error) {
	var todos []*Todo
	_, err := db.Select(&todos, "SELECT * FROM " + config.TABLE_NAME_TODOS + " ORDER BY " + config.TABLE_NAME_TODOS + ".update DESC")
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (todo *Todo) IsEnable() bool {
    return todo.Active != -1
}

func (todo *Todo) PreInsert(s gorp.SqlExecutor) error {
	todo.CreateAt = time.Now().Unix()
	todo.UpdateAt = todo.CreateAt
	todo.ProcCount = 1
	todo.Status = 0
	todo.Active = 1
	return nil
}

func (todo *Todo) PreUpdate(s gorp.SqlExecutor) error {
	todo.UpdateAt = time.Now().Unix()
	return nil
}

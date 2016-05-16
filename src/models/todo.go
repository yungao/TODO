package models

import (
	"github.com/coopernurse/gorp"
	"log"
	"time"
)

type Todo struct {
	ID       int    `db:"id"        json:"id"`
	Name     string `db:"name"      json:"name"        form:"name"     binding:"required"`
	Title    string `db:"title"     json:"title"       form:"title"    binding:"required"`
	Content  string `db:"content"   json:"content"     form:"content"`
	Priority int8   `db:"priority"  json:"priority"    form:"priority" binding:"required"`
	/* Type:
	 *       0:  public
	 *       1:  private
	 */
	Type     int8  `db:"type"      json:"type"        form:"type"     binding:"required"`
	CreateAt int64 `db:"create"    json:"create"`
	LimitAt  int64 `db:"limit"     json:"limit"       form:"limit"    binding:"required"`
	UpdateAt int64 `db:"update"    json:"update"`
	/* Status:
	 *       0: undone
	 *       1: done
	 */
	Status    int8 `db:"status"    json:"status"`
	ProcCount int  `db:"count"     json:"count"`
	/* Active:
	 *       -1: deleted
	 *       0:  diabled
	 *       1:  normal
	 */
	Active int8 `db:"active"    json:"active"`
}

type Todos struct {
	Collection []Todo `json:"todos"`
}

// Create todo table if not exist
func CreateTodoTable(db *gorp.DbMap) {
	tb := db.AddTableWithName(Todo{}, "tb_todos")
	tb.SetKeys(true, "id")
	tb.ColMap("name").SetMaxSize(100).SetNotNull(true)
	tb.ColMap("title").SetMaxSize(255).SetNotNull(true)
	tb.ColMap("content").SetMaxSize(2048)
	tb.ColMap("priority").SetNotNull(true)
	tb.ColMap("type").SetNotNull(true)
	tb.ColMap("create").SetNotNull(true)
	tb.ColMap("limit").SetNotNull(true)
	tb.ColMap("update").SetNotNull(true)
	tb.ColMap("status").SetNotNull(true)
	tb.ColMap("count").SetNotNull(true)
	tb.ColMap("active").SetNotNull(true)

	err := db.CreateTablesIfNotExists()
	// db.DropTables()
	if err != nil {
		panic(err)
	}

	log.Println(">>> Todo table created")
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

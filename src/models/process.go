package models

import (
	"config"
	"github.com/coopernurse/gorp"
	"log"
	"time"
	"fmt"
)

type Process struct {
	ID        int    `db:"id"        json:"id"`
	TodoID    int    `db:"todoid"    json:"todoid"    form:"todoid"      binding:"required"`
	CreatorID int    `db:"uid"       json:"-"`
	Creator   *User  `db:"-"         json:"creator"`
	Partners  string `db:"partners"  json:"partners"  form:"partners"`
	Tags      string `db:"tags"      json:"tags"      form:"tags"`
	/** Action:
	 *       0:  Complete Todo
	 *       1:  Create Todo
	 *       2:  Add partner
	 *       3:  Add Tag
	 *       4:  Update Todo
	 *       5:  Assigned to others
	 *       6:  Revoked Todo
	 */
	Action      int8   `db:"action"    json:"action"    form:"action"   binding:"required"`
	Content     string `db:"content"   json:"content"   form:"content"`
	AttachCount int    `db:"fcount"    json:"fcount"`
	UpdateAt    int64  `db:"update"    json:"update"`
	UserAgent   string `db:"agent"     json:"agent"`
	/* Active:
	 *       -1: deleted
	 *       1:  normal
	 */
	Active int8 `db:"active"    json:"active"`
}

type Processes struct {
	Collection []Process `json:"processes"`
}

// Create process table if not exist
func CreateProcessTable(db *gorp.DbMap) {
	tb := db.AddTableWithName(Process{}, config.TABLE_NAME_TODO_PROCESS)
	tb.SetKeys(true, "id")
	tb.ColMap("todoid").SetNotNull(true)
	tb.ColMap("uid").SetNotNull(true)
	// tb.ColMap("asgid")
	// tb.ColMap("tagid")
	tb.ColMap("action").SetNotNull(true)
	tb.ColMap("content").SetMaxSize(2048)
	tb.ColMap("fcount").SetNotNull(true)
	tb.ColMap("update").SetNotNull(true)
	tb.ColMap("agent").SetMaxSize(100)
	tb.ColMap("active").SetNotNull(true)

	err := db.CreateTablesIfNotExists()
	// db.DropTables()
	if err != nil {
		panic(err)
	}

	log.Printf(">>> Table[%s] created", config.TABLE_NAME_TODO_PROCESS)
}

/*
* Query todo processes from database by todo ID
 */
func GetTodoProcesses(db *gorp.DbMap, id int) ([]*Process, error) {
	var processes []*Process
	_, err := db.Select(&processes, fmt.Sprintf("SELECT * FROM %s  WHERE todoid=%d ORDER BY %s.update DESC", config.TABLE_NAME_TODO_PROCESS, id, config.TABLE_NAME_TODO_PROCESS))
	if err != nil {
		return nil, err
	}

	return processes, nil
}

func (proc *Process) PreInsert(s gorp.SqlExecutor) error {
	proc.UpdateAt = time.Now().Unix()
	proc.Active = 1
	return nil
}

// func (proc *Process) PreUpdate(s gorp.SqlExecutor) error {
// 	proc.UpdateAt = time.Now().Unix()
// 	return nil
// }


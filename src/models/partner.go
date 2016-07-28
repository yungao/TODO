package models

import (
	"config"
	"github.com/coopernurse/gorp"
	"log"
	"fmt"
)

type Partner struct {
	ID        int  `db:"id"        json:"id"`
	TodoID    int  `db:"todoid"    json:"todoid"   form:"todoid"   binding:"required"`
	PartnerID int  `db:"pid"       json:"pid"      form:"pid"      binding:"required"`
	Partner   *User `db:"-"         json:"partner"`
	UserID    int  `db:"uid"       json:"uid"`
	Creator   *User `db:"-"         json:"creator"`
	Duty      int8 `db:"duty"      json:"duty"     form:"duty"`
	/* Active:
	 *       -1: deleted
	 *       1:  normal
	 */
	Active int8 `db:"active"    json:"active"`
}

type Partners struct {
	Collection []*Partner `json:"partners"`
}

// Create partner table if not exist
func CreatePartnerTable(db *gorp.DbMap) {
	tb := db.AddTableWithName(Partner{}, config.TABLE_NAME_TODO_PARTNER)
	tb.SetKeys(true, "id")
	tb.ColMap("todoid").SetNotNull(true)
	tb.ColMap("pid").SetNotNull(true)
	tb.SetUniqueTogether("todoid", "pid")
	tb.ColMap("uid").SetNotNull(true)
	tb.ColMap("duty").SetNotNull(true)
	tb.ColMap("active").SetNotNull(true)

	err := db.CreateTablesIfNotExists()
	// db.DropTables()
	if err != nil {
		panic(err)
	}

	log.Printf(">>> Table[%s] created", config.TABLE_NAME_TODO_PARTNER)
}

/*
* Query todo partners from database by TodoID
 */
func GetPartnersByTodoID(db *gorp.DbMap, tid int) ([]*Partner, error) {
	var partners []*Partner
	_, err := db.Select(&partners, fmt.Sprintf("SELECT * FROM %s WHERE %s.todoid = %d", config.TABLE_NAME_TODO_PARTNER, config.TABLE_NAME_TODO_PARTNER, tid))
	if err != nil {
		return nil, err
	}

	return partners, nil
}

func (partner *Partner) PreInsert(s gorp.SqlExecutor) error {
	partner.Active = 1
	return nil
}

func (partner *Partner) String() string {
    return fmt.Sprintf("{ID:%d, TodoID:%d, PartnerID:%s, Duty:%s, Active:%d}", partner.ID, partner.TodoID, partner.UserID, partner.Duty, partner.Active)
}

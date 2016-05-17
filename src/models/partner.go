package models

import (
	"config"
	"github.com/coopernurse/gorp"
	"log"
)

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

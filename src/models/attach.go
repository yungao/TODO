package models

import (
	"config"
	"github.com/coopernurse/gorp"
	"log"
	"time"
)

type Attach struct {
	ID        int    `db:"id"        json:"id"`
	CreatorID int    `db:"uid"       json:"uid"`
	TodoID    int    `db:"todoid"    json:"todoid"      form:"todoid"   binding:"required`
	ProcessID int    `db:"procid"    json:"procid"      form:"procid"   binding:"required`
	Name      string `db:"name"      json:"name"        form:"name"     binding:"required"`
	Path      string `db:"path"      json:"path"        form:"path"     binding:"required"`
	Type      string `db:"type"      json:"type"        form:"type"     binding:"required"`
	Remark    string `db:"remark"    json:"remark"      form:"remark"`
	UploadAt  int64  `db:"upload"    json:"upload"`
	/* Active:
	 *       -1: deleted
	 *       -1: disabled
	 *       1:  normal
	 */
	Active int8 `db:"active"    json:"active"`
}

type Attachs struct {
	Collection []Attach `json:"attachs"`
}

// Create attach table if not exist
func CreateAttachTable(db *gorp.DbMap) {
	tb := db.AddTableWithName(Attach{}, config.TABLE_NAME_TODO_ATTACHS)
	tb.SetKeys(true, "id")
	tb.ColMap("uid").SetNotNull(true)
	tb.ColMap("todoid").SetNotNull(true)
	tb.ColMap("procid").SetNotNull(true)
	tb.ColMap("name").SetMaxSize(255).SetNotNull(true)
	tb.ColMap("path").SetMaxSize(4096).SetNotNull(true)
	tb.ColMap("type").SetMaxSize(10).SetNotNull(true)
	tb.ColMap("remark").SetMaxSize(255)
	tb.ColMap("upload").SetNotNull(true)
	tb.ColMap("active").SetNotNull(true)

	err := db.CreateTablesIfNotExists()
	// db.DropTables()
	if err != nil {
		panic(err)
	}

	log.Printf(">>> Table[%s] created", config.TABLE_NAME_TODO_ATTACHS)
}

func (attach *Attach) PreInsert(s gorp.SqlExecutor) error {
	attach.UploadAt = time.Now().Unix()
	attach.Active = 1
	return nil
}

// func (tag *Tag) PreUpdate(s gorp.SqlExecutor) error {
// 	tag.UpdateAt = time.Now().Unix()
// 	return nil
// }

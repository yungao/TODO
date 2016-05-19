package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
)

const (
	TABLE_NAME_USER         = "tb_users"
	TABLE_NAME_TAGS         = "tb_tags"
	TABLE_NAME_TODOS        = "tb_todos"
	TABLE_NAME_TODO_PARTNER = "tb_todo_partners"
	TABLE_NAME_TODO_PROCESS = "tb_todo_process"
	TABLE_NAME_TODO_TAGS    = "tb_todo_tags"
	TABLE_NAME_TODO_ATTACHS = "tb_todo_attachs"
)

func EnableDBLogger(db *gorp.DbMap) {
	db.TraceOn("[gorp]", log.New(os.Stdout, "[ MySQL ] ", log.Lmicroseconds))
	//dbp.TraceOff()
}

//var dbMap *gorp.DbMap
//
//func GetDBMap() *gorp.DbMap {
//    if dbMap == nil {
//        DB()
//    }
//
//    return dbMap
//}

func DB() *gorp.DbMap {
	conn, err := sql.Open("mysql", "root:@/gotodo")

	if err != nil {
		panic(err)
	}

	dbMap := &gorp.DbMap{Db: conn, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	return dbMap
}

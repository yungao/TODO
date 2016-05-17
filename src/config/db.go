package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"

	model "models"
)

func initDBTables(db *gorp.DbMap) {
	model.CreateUserTable(db)
	model.CreateTodoTable(db)
	model.CreateTagTable(db)
	model.CreateTodoTagTable(db)
	model.CreateProcessTable(db)
}

func enableDBLogger(db *gorp.DbMap) {
	db.TraceOn("[gorp]", log.New(os.Stdout, ">>MySQL<< ", log.Lmicroseconds))
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

	initDBTables(dbMap)
	enableDBLogger(dbMap)

	return dbMap
}

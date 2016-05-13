package config

import (
    "log"
    "os"
    "database/sql"

    "github.com/coopernurse/gorp"
    _ "github.com/go-sql-driver/mysql"

    model   "models"
)

func initDBTables(db *gorp.DbMap) {
    // Add [user] table
    tb := db.AddTableWithName(model.User{}, "user").SetKeys(true, "ID")
    tb.ColMap("name").SetUnique(true)

    db.CreateTables()
    //db.DropTables()
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
    conn, err := sql.Open("mysql", "root:@/goweb")

    if err != nil {
        panic(err)
    }

    dbMap := &gorp.DbMap{Db: conn, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

    initDBTables(dbMap);
    enableDBLogger(dbMap);

    return dbMap
}

package config

import (
    "log"
    "os"
    "database/sql"

    "github.com/coopernurse/gorp"
    _ "github.com/go-sql-driver/mysql"

    "model"
)

func initDBTables(db *gorp.DbMap) {
    db.AddTableWithName(model.User{}, "user").SetKeys(true, "ID")
    db.CreateTables()
    //db.DropTables()
}

func enableDBLogger(db *gorp.DbMap) {
    db.TraceOn("[gorp]", log.New(os.Stdout, ">>MySQL: ", log.Lmicroseconds))
    //dbp.TraceOff()
}

func DB() *gorp.DbMap {
    conn, err := sql.Open("mysql", "root:@/goweb")

    if err != nil {
        panic(err)
    }

    dbmap := &gorp.DbMap{Db: conn, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

    initDBTables(dbmap);
    enableDBLogger(dbmap);

    return dbmap
}

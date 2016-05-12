package config

import (
    //"log"
    //"os"
    "database/sql"

    "github.com/coopernurse/gorp"
    _ "github.com/go-sql-driver/mysql"
)

var dbMap *gorp.DbMap

func GetDBMap() *gorp.DbMap {
    if dbMap == nil {
        DB()
    }

    return dbMap
}

func DB() *gorp.DbMap {
    conn, err := sql.Open("mysql", "root:@/goweb")

    if err != nil {
        panic(err)
    }

    dbMap = &gorp.DbMap{Db: conn, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

    return dbMap
}

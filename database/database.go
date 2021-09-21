package database

import (
    "database/sql"
    "log"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

const Source string = "./fl.db"

var DB *sql.DB

func init() {
    db, err := sql.Open("sqlite3", Source)
    if err != nil {
        log.Fatal(err)
    }
    db.SetMaxOpenConns(60)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(30 * time.Minute)
    DB = db
}

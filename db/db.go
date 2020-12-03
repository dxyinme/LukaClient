package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var databaseConnector *sql.DB

func NewConn(dbFile string) (err error) {
	databaseConnector , err = sql.Open("sqlite3", dbFile)
	return
}

func ExecCmd(cmd string) (sql.Result, error) {
	return databaseConnector.Exec(cmd)
}

func Close() (err error) {
	err = databaseConnector .Close()
	return
}
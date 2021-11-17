package main

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"github.com/povsister/scp"
)

func main() {

	db, _ := sql.Open("sqlite3", "./filtros.db")
	stmt, _ := db.Prepare("CREATE TABLE IF NOT EXISTS 'test' { 'ID' INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 'content' TEXT }")
	stmt.Exec()

}
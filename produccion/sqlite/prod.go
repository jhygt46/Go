package main

import (
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	//"github.com/povsister/scp"
)

func main() {

	db, err1 := sql.Open("sqlite3", "./filtros.db")
	if err1 != nil {
		fmt.Println(err1)
	}
	stmt, err2 := db.Prepare("CREATE TABLE IF NOT EXISTS 'test' { 'ID' INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 'content' TEXT }")
	if err2 != nil {
		fmt.Println(err2)
	}
	stmt.Exec()

}
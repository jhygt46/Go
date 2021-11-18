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
	fmt.Printf("%T", db)
	//create_db(&db)

}
/*
func create_db(*db ){
	stmt, err2 := *db.Prepare(`create table if not exists user (id  integer not null primary key,name text, age integer)`)
	if err2 != nil {
		fmt.Println(err2)
	}
	stmt.Exec()
}
*/
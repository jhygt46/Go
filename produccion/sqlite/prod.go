package main

import (
	"os"
	"fmt"
	"time"
	//"context"
	"strconv"
	"math/big"
	"io/ioutil"
	"crypto/rand"
	"database/sql"
	"path/filepath"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	//"github.com/povsister/scp"
)

func main() {

	now := time.Now()
	db, err := sql.Open("sqlite3", "./filtros1.db")
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
	}
	printelaped(now, "OPEN DB")
	create_db(db)

	/*
	stmt, err := db.PrepareContext(ctx, "SELECT content FROM contents WHERE id=?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	id := 43
	var username string
	err = stmt.QueryRowContext(ctx, id).Scan(&username)
	fmt.Println(err)
	fmt.Println(username)
	*/

	//create_db(db)
	/*
	escribir_db(db, 10000)
	select_db(db, 2500, 160000)
	escribir_db(db, 10000)
	select_db(db, 2500, 170000)
	escribir_db(db, 10000)
	select_db(db, 2500, 180000)
	escribir_db(db, 10000)
	select_db(db, 2500, 190000)
	escribir_db(db, 10000)
	select_db(db, 2500, 200000)
	escribir_db(db, 10000)
	select_db(db, 2500, 210000)
	escribir_db(db, 10000)
	select_db(db, 2500, 220000)
	escribir_db(db, 10000)
	select_db(db, 2500, 230000)
	escribir_db(db, 10000)
	select_db(db, 2500, 240000)
	escribir_db(db, 10000)
	select_db(db, 2500, 250000)
	*/
	//escribir_file("/var/db1_test", 25000)
	//select_file("/var/db1_test", 25000)

}

type Objecto struct {
	Name string `json:"Name"`
}

func select_db(db *sql.DB, numb int, max int){
	c := 0
	now := time.Now()
	for n := 0; n < numb; n++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
		content := get_content(db, n.Int64());
		if content == ""{
			fmt.Println("CONTENT VACIO")
		}
		//readcon(content)
		c++
	}
	elapsed := time.Since(now)
	fmt.Printf("SELECT %v [%s] c/u\n", c, time_cu(elapsed, c))
}
func select_file(path string, numb int){
	c := 0
	//numb := 250000
	now := time.Now()
	for n := 0; n < numb; n++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(2000))
		folder := getFolder64(n.Uint64())
		file, err := os.Open(path+"/"+folder+"/56")
		if err != nil{
			fmt.Println(err)
		}
		file.Close()
		byteValue, _ := ioutil.ReadAll(file)
		read(byteValue)
		c++
	}
	elapsed := time.Since(now)
	fmt.Printf("SELECT %v FILES en [%v]\n", c, elapsed)
}
func escribir_file(path string, numb int){

	d1 := []byte("{\"Id\":1,\"Data\":{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}}")
	c := 0

	numb = numb / 100
	now := time.Now()
	for n := 0; n < numb; n++ {
		folder := getFolder64(uint64(n*100))
		newpath := filepath.Join(path, folder)
		err := os.MkdirAll(newpath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			fmt.Println("FOLDER ERROR: ", err)
		}
		for i := 0; i < 100; i++ {
			err := os.WriteFile(path+"/"+folder+"/"+strconv.Itoa(i), d1, 0644)
			if err != nil {
				fmt.Println(err)
			}
			c++
		}
	}
	elapsed := time.Since(now)
	fmt.Printf("WRITE FILES %v en [%v]\n", c, elapsed)
}
func escribir_db(db *sql.DB, numb int){
	d1 := []byte("{\"Id\":1,\"Data\":{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}}")
	now := time.Now()
	//now1 := time.Now()
	c := 0
	for n := 0; n < numb; n++ {
		add_txt_db(db, string(d1))
		c++
		/*
		if c % 1000 == 0 { 
			elapsed1 := time.Since(now1)
			fmt.Printf("WRITE 1000 [%s] c/u\n", time_cu(elapsed1, c))
			now1 = time.Now()
		}
		*/
	}
	elapsed := time.Since(now)
	fmt.Printf("WRITE DB %v en [%v] [%s] c/u\n", c, elapsed, time_cu(elapsed, c))
}
func create_db(db *sql.DB){
	now := time.Now()
	stmt, err := db.Prepare(`create table if not exists contents (id integer not null primary key,content text)`)
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()
	printelaped(now, "CREATE DB")
}
func add_obj_db(db *sql.DB, obj Objecto){
	stmt, err := db.Prepare("INSERT INTO contents(content) values(?)")
	if err != nil {
		fmt.Println(err)
	}
	u, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	res, err := stmt.Exec(u)
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
	}
}
func add_txt_db(db *sql.DB, str string){
	stmt, err := db.Prepare("INSERT INTO contents(content) values(?)")
	if err != nil {
		fmt.Println(err)
	}
	res, err := stmt.Exec(str)
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
	}
}
func update_db(db *sql.DB, id int64){
	stmt, err := db.Prepare("UPDATE contents SET content=? WHERE id=?")
	if err != nil {
		fmt.Println(err)
	}
	res, err := stmt.Exec("content", id)
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
	}
}
func get_content(db *sql.DB, id int64) string {
	rows, err := db.Query("SELECT content FROM contents WHERE id=?", id)
	if err != nil { panic(err) }
	defer rows.Close()
	var content string
	for rows.Next() {
		err2 := rows.Scan(&content)
		if err2 != nil { fmt.Println(err2) }
	}
	return content
}
func getFolder64(num uint64) string {

	c1, n1 := divmod(num, 1000000)
	c2, n2 := divmod(n1, 10000)
	c3, _ := divmod(n2, 100)
	return strconv.FormatUint(c1, 10)+"/"+strconv.FormatUint(c2, 10)+"/"+strconv.FormatUint(c3, 10)

}
func divmod(numerator, denominator uint64) (quotient, remainder uint64) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}
func printelaped(start time.Time, str string) /*time.Time*/ {
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
	//return time.Now()
}
func read(x []byte){
	//
}
func readcon(x string){
	//
}
func time_cu(t time.Duration, c int) string {
	ms := float64(t / time.Nanosecond)
	res := ms / float64(c)
	var s string
	if res < 1000 {
		s = fmt.Sprintf("%.2f NanoSec", res)
	} else if res >= 1000 && res < 1000000{
		s = fmt.Sprintf("%.2f MicroSec", res/1000)
	} else {
		s = fmt.Sprintf("%.2f MilliSec", res/1000000)
	}
	return s
}
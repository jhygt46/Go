package main

import (
	"fmt"
	"time"
	"crypto/rand"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	//"github.com/povsister/scp"
)

func main() {

	now := time.Now()
	db, err1 := sql.Open("sqlite3", "./filtros1.db")
	if err1 != nil {
		fmt.Println(err1)
	}
	now = printelaped(now, "OPEN DB")
	create_db(db)
	now = printelaped(now, "CREATE DB")
	select_db(db)
	now = printelaped(now, "SELECT DB")
	escribir_db(db, "PRUEBA")
	now = printelaped(now, "WRITE DB")

	select_file()
	now = printelaped(now, "SELECT FILES")
	escribir_file("/var/db1_test/")
	now = printelaped(now, "WRITE FILES")

}

func select_db(db *sql.DB){
	numb := 100000
	now := time.Now()
	for n := 0; n < numb; n++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(100000))
		content := get_content(&db, n.Uint64());
	}
	printelaped(now, "SELECT DB")
}
func select_file(){
	numb := 100000
	now = time.Now()
	for n := 0; n < numb; n++ {
		n, _ := remove.Int(remove.Reader, big.NewInt(1000000))
		folder := getFolder64(n.Uint64())
		file, err := os.Open(path+"/"+folder)
		if err != nil{
			fmt.Println(err)
		}
		file.Close()
		byteValue, _ := ioutil.ReadAll(file)
		read(byteValue)
	}
	printelaped(now, "SELECT FILES")
}
func escribir_file(path string){

	d1 := []byte("{\"Id\":1,\"Data\":{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}}")
	c := 0

	numb := 800
	now = time.Now()
	for n := 0; n < numb; n++ {

		v := 100
		folder := getFolder64(uint64(n))

		newpath := filepath.Join(path, folder)
		err := os.MkdirAll(newpath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			fmt.Println("FOLDER ERROR: ", err)
		}

		for i := 0; i < v; i++ {
			err := os.WriteFile(path+"/"+folder+"/"+strconv.Itoa(i), d1, 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
		c++

	}
	elapsed := time.Since(now)
	fmt.Printf("Cantidad %v / Tiempo: [%v]\n", c, elapsed)
}
func escribir_db(db *sql.DB, str string){
	numb := 80000
	now = time.Now()
	c := 0
	for n := 0; n < numb; n++ {
		add_txt_db(&db, strP)
		c++
	}
	elapsed := time.Since(now)
	fmt.Printf("Cantidad %v / Tiempo: [%v]\n", c, elapsed)
}

func create_db(db *sql.DB){
	now = time.Now()
	stmt, err := db.Prepare(`create table if not exists contents (id integer not null primary key,content text)`)
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()
	printelaped(now, "CREATE DB")
}
func add_obj_db(db *sql.DB, obj objeto){
	stmt, err := db.Prepare("INSERT INTO contents(content) values(?)")
	if err != nil {
		fmt.Println(err)
	}
	u, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	_, err := stmt.Exec(u)
	if err != nil {
		fmt.Println(err)
	}
}
func add_txt_db(db *sql.DB, str string){
	stmt, err := db.Prepare("INSERT INTO contents(content) values(?)")
	if err != nil {
		fmt.Println(err)
	}
	_, err := stmt.Exec(str)
	if err != nil {
		fmt.Println(err)
	}
}
func update_db(db *sql.DB, id int64){
	stmt, err := db.Prepare("UPDATE contents SET content=? WHERE id=?")
	if err != nil {
		fmt.Println(err)
	}
	_, err := stmt.Exec("content", id)
	if err != nil {
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
	}
	return content
}
func get_contents(db *sql.DB, id int64) string {
	rows, err := db.Query("SELECT content FROM contents WHERE id=?", id)
	if err != nil { panic(err) }
	defer rows.Close()
	var content string
	for rows.Next() {
		err2 := rows.Scan(&content)
	}
	return content
}

func printelaped(start time.Time, str string) time.Time {
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
	return time.Now()
}
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
	"github.com/valyala/fasthttp"
	//"github.com/povsister/scp"
)
type Config struct {

}
type MyHandler struct {
	Dbs []*sql.DB `json:"Dbs"`
	Config Config `json:"Config"`
	Minicache map[uint64]*Data
}
type Data struct {
	C [] Campos `json:"C"`
	E [] Evals `json:"E"`
}
type Campos struct {
	T int `json:"T"`
	N string `json:"N"`
	V [] string `json:"V"`
}
type Evals struct {
	T int `json:"T"`
	N string `json:"N"`
}
func main() {

	len := 10

	dbs := make([]*sql.DB, 0)
	now := time.Now()
	for i:=0; i < len; i++ {
		db, err := sql.Open("sqlite3", "./filtros"+strconv.Itoa(i)+".db")
		if err == nil {
			stmt, err := db.Prepare(`create table if not exists contents (id integer not null primary key autoincrement,content text)`)
			if err != nil {
				fmt.Println(err)
			}
			stmt.Exec()
			dbs = append(dbs, db)
		}else{
			fmt.Println(err)
		}
	}
	printelaped(now, "OPEN DBS")

	h := &MyHandler{ Dbs: dbs }
	fasthttp.ListenAndServe(":81", h.HandleFastHTTP)

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


}

type Objecto struct {
	Name string `json:"Name"`
}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	db := read_int64(ctx.QueryArgs().Peek("db"))
	id := read_int64(ctx.QueryArgs().Peek("id")) 
	total := read_int64(ctx.QueryArgs().Peek("total")) 

	switch string(ctx.Path()) {
	case "/get-op1":
		select_db(h.Dbs[db], 2500, total)
		fmt.Fprintf(ctx, "OK")
	case "/get-op2":
		select_db_rand(h.Dbs, 2500, total)
		fmt.Fprintf(ctx, "OK")
	case "/get-op3":
		select_file("/var/db1_test", 2500, total)
		fmt.Fprintf(ctx, "OK")
	case "/get-op4":
		h.select_memory(2500, total)
		fmt.Fprintf(ctx, "OK")
	case "/put-op1":
		escribir_db(h.Dbs[db], id)
		fmt.Fprintf(ctx, "OK")
	case "/put-op2":
		escribir_file("/var/db1_test", total)
		fmt.Fprintf(ctx, "OK")
	case "/put-op3":
		h.escribir_memory("/var/db1_test", int(total))
		fmt.Fprintf(ctx, "OK")
	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}
	
}
func (h *MyHandler) escribir_memory(path string, numb int){

	c := 0
	now := time.Now()
	for n := 1; n <= numb; n++ {
		jsonFiltro, err := os.Open(path+"/"+getFolderFile64(int64(n)))
		if err == nil {
			byteValueFiltro, _ := ioutil.ReadAll(jsonFiltro)
			data := Data{}
			if err := json.Unmarshal(byteValueFiltro, &data); err == nil {
				h.Minicache[uint64(n)] = &data
			}
		}
		defer jsonFiltro.Close()
	}
	elapsed := time.Since(now)
	fmt.Printf("SELECT %v [%s] c/u total %v\n", c, time_cu(elapsed, c), elapsed)

}
func (h *MyHandler) select_memory(numb int, max int64){

	c := 0
	now := time.Now()
	for n := 0; n < numb; n++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(max))
		if res, found := h.Minicache[n.Uint64()]; found {
			readdata(res)
			c++
		}
	}
	elapsed := time.Since(now)
	fmt.Printf("SELECT %v [%s] c/u total %v\n", c, time_cu(elapsed, c), elapsed)

}
func select_db_rand(db []*sql.DB, numb int, max int64){
	c := 0
	now := time.Now()
	for n := 0; n < numb; n++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(max))
		r, _ := rand.Int(rand.Reader, big.NewInt(10))
		content := get_content(db[r.Int64()], n.Int64());
		if content == ""{
			fmt.Println("CONTENT VACIO")
		}
		//readcon(content)
		c++
	}
	elapsed := time.Since(now)
	fmt.Printf("SELECT %v [%s] c/u total %v\n", c, time_cu(elapsed, c), elapsed)
}
func select_db(db *sql.DB, numb int, max int64){
	c := 0
	now := time.Now()
	for n := 0; n < numb; n++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(max))
		content := get_content(db, n.Int64());
		if content == ""{
			fmt.Println("CONTENT VACIO")
		}
		//readcon(content)
		c++
	}
	elapsed := time.Since(now)
	fmt.Printf("SELECT %v [%s] c/u total %v\n", c, time_cu(elapsed, c), elapsed)
}
func select_file(path string, numb int, max int64){
	c := 0
	now := time.Now()
	for n := 0; n < numb; n++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(max))
		folder := getFolder64(n.Int64())
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
	fmt.Printf("SELECTFILES %v [%s] c/u total %v\n", c, time_cu(elapsed, c), elapsed)
}
func escribir_file(path string, numb int64){

	d1 := []byte("{\"Id\":1,\"Data\":{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}}")
	c := 0

	numb2 := int(numb / 100)
	now := time.Now()
	for n := 0; n < numb2; n++ {
		folder := getFolder64(int64(n*100))
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
func escribir_db(db *sql.DB, numb int64){
	d1 := []byte("{\"Id\":1,\"Data\":{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}}")
	now := time.Now()
	//now1 := time.Now()
	c := 0
	for n := 0; n < int(numb); n++ {
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
func getFolder64(num int64) string {

	c1, n1 := divmod(num, 1000000)
	c2, n2 := divmod(n1, 10000)
	c3, _ := divmod(n2, 100)
	return strconv.FormatInt(c1, 10)+"/"+strconv.FormatInt(c2, 10)+"/"+strconv.FormatInt(c3, 10)

}
func getFolderFile64(num int64) string {

	c1, n1 := divmod(num, 1000000)
	c2, n2 := divmod(n1, 10000)
	c3, c4 := divmod(n2, 100)
	return strconv.FormatInt(c1, 10)+"/"+strconv.FormatInt(c2, 10)+"/"+strconv.FormatInt(c3, 10)+"/"+strconv.FormatInt(c4, 10)

}
func divmod(numerator, denominator int64) (quotient, remainder int64) {
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
func readdata(x *Data){
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
func read_int64(data []byte) int64 {
    var x int64
    for _, c := range data {
        x = x * 10 + int64(c - '0')
    }
    return x
}
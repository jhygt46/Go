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
type Objecto struct {
	Name string `json:"Name"`
}
type Data struct {
	C [] Campos `json:"C"`
	E [] Evals `json:"E"`
	N string `json:"N"`
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
type MyHandler struct {
	Dbs *sql.DB `json:"Dbs"`
	Config Config `json:"Config"`
	//Minicache map[int64]*Data `json:"Minicache"`
	Minicache *Minicache `json:"Minicache"`
	DB []Dbs `json:"DB"`
	MDBS []*sql.DB `json:"MDBS"`
}
type Minicache struct {
	Cache map[int64]*Data `json:"Cache"`
}
type Dbs struct {
	Db *sql.DB `json:"Db"`
	Stmt []*sql.Stmt `json:"Stmt"`
}
//context.Background()

var cache map[int64]*Data

func main() {

	/*
	dbs := make([]Dbs, 0)
	len := 10
	for i:=0; i<len; i++ {
		db, err := getsqlite(i)
		if err == nil {
			stmt, err := db.Prepare("SELECT content FROM contents WHERE id=?")
			if err == nil {
				arrstmt := make([]*sql.Stmt, 0)
				arrstmt = append(arrstmt, stmt)
				dbs = append(dbs, Dbs{ Db: db, Stmt: arrstmt })
			}
		}
	}
	*/

	for i:=1; i<=350000; i++ {
		cache[int64(i)] = &Data{}
	}

	dbs := make([]*sql.DB, 0)
	len := 10
	for i:=0; i<len; i++ {
		db, err := getsqlite(i)
		if err == nil {
			dbs = append(dbs, db)
		}
	}


	h := &MyHandler{ MDBS: dbs }
	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)

	//escribir_file("/var/db1_test", 3500)

	//total := 200000
	//db, err := getsqlite(0)
	/*
	if err == nil {
		cache := make(map[int64]*Data, total)
		now := time.Now()
		for i:=1; i<=total; i++ {
			folderfile := getFolderFile64(random(int64(i)))
			file, err := os.Open("/var/db1_test/"+folderfile)
			if err != nil{
				fmt.Println(err)
			}
			byteValue, err := ioutil.ReadAll(file)
			file.Close()
			data := Data{}
			if err := json.Unmarshal(byteValue, &data); err == nil {
				cache[int64(i)] = &data
			}
		}
		printelaped(now, "CACHE LISTO")
		minicache := &Minicache{ Cache: cache }
		h := &MyHandler{ MDBS: dbs }
		fasthttp.ListenAndServe(":80", h.HandleFastHTTP)	
	}
	*/
	
}
func getsqlite(i int) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./filtros"+strconv.Itoa(i)+".db")
	if err == nil {
		stmt, err := db.Prepare(`create table if not exists contents (id integer not null primary key autoincrement,content text)`)
		if err != nil {
			fmt.Println("err1")
			fmt.Println(err)
			return db, err
		}
		stmt.Exec()
		return db, nil
	}else{
		fmt.Println("err2")
		fmt.Println(err)
		return db, err
	}
}
func get_content(db *sql.DB, id int64) (string, error) {
	rows, err := db.Query("SELECT content FROM contents WHERE id=?", id)
	if err != nil { 
		return "", err
	}
	defer rows.Close()
	var content string
	for rows.Next() {
		err := rows.Scan(&content)
		if err != nil { 
			return "", err
		}
	}
	return content, nil
}
func get_content2(db *sql.DB, id int64) (string, error) {
	stmt, err := db.Prepare("SELECT content FROM contents WHERE id=?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var content string
	errs := stmt.QueryRow(id).Scan(&content)
	if errs != nil {
		if errs == sql.ErrNoRows {
			// Handle the case of no rows returned.
		}
		return "", errs
	}else{
		return content, nil
	}
}
func add_txt_db(db *sql.DB) (error) {

	str := []byte("{\"Id\":1,\"Data\":{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}}")
	stmt, err := db.Prepare("INSERT INTO contents (content) values (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(string(str))
	if err != nil {
		return err
		fmt.Println(res)
	}
	return nil

}
func escribir_file(path string, numb int){

	d1 := []byte("{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}")
	c := 0

	now := time.Now()
	for n := 0; n < numb; n++ {
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
	fmt.Printf("WRITES FILES %v [%s] c/u total %v\n", c, time_cu(elapsed, c), elapsed)

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	

	ctx.Response.Header.Set("Content-Type", "application/json")
	
	switch string(ctx.Path()) {
	case "/get0":
		
		x := random(300000)
		if res, found := h.Minicache.Cache[x]; found {
			json.NewEncoder(ctx).Encode(res)
		}else{
			content, err := get_content(h.Dbs, x)
			if err == nil{
				fmt.Fprintf(ctx, content)
			}else{
				ctx.Error("Not Found", fasthttp.StatusNotFound)
			}
		}

	case "/get1":
		
		db, id := getdbid(random(1000000), 0)
		content, err := get_content(h.MDBS[db], id)
		if err == nil{
			fmt.Fprintf(ctx, content)
		}else{
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
		
	case "/get2":
		
		folderfile := getFolderFile64(random(300000))
		file, err := os.Open("/var/db1_test/"+folderfile)
		if err != nil{
			fmt.Println(err)
		}
		byteValue, err := ioutil.ReadAll(file)
		if err != nil{
			fmt.Println(err)
		}
		file.Close()
		fmt.Fprintf(ctx, string(byteValue))

	case "/put1":
		
		id := read_int64(ctx.QueryArgs().Peek("id"))

		str1 := []byte("{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}")
		str := string(str1)
		tx, err := h.MDBS[id].Begin()
		if err != nil {
			fmt.Println(err)
		}
		defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.
		stmt, err := tx.Prepare("INSERT INTO contents(content) VALUES(?)")
		if err != nil {
			fmt.Println(err)
		}
		defer stmt.Close() // Prepared statements take up server resources and should be closed after use.
		now := time.Now()
		for i:=0; i<100000; i++ {
			if _, err := stmt.Exec(str); err != nil {
				fmt.Println(err)
			}
		}
		printelaped(now, "INSERT 100000")
		if err := tx.Commit(); err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(ctx, "OK")

	case "/update":
		
		//str1 := []byte("{\"Id\":1,\"Data\":{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Alemana\", \"Española\", \"Inglesa\", \"Francesa\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}}")
		str1 := []byte("{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Alemana\", \"Española\", \"Inglesa\", \"Francesa\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}")
		str := string(str1)
		tx, err := h.Dbs.Begin()
		if err != nil {
			fmt.Println(err)
		}
		defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.
		stmt, err := tx.Prepare("UPDATE contents SET content=? WHERE id=?")
		if err != nil {
			fmt.Println(err)
		}
		defer stmt.Close() // Prepared statements take up server resources and should be closed after use.
		now := time.Now()
		for i:=1; i<=300000; i++ {
			if _, err := stmt.Exec(str, i); err != nil {
				fmt.Println(err)
			}
		}
		printelaped(now, "UPDATE 300000")
		if err := tx.Commit(); err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(ctx, "OK")

	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}

	/*


	switch string(ctx.Path()) {
	case "/get-op1":
		select_db(h.Dbs[db], int(times), total)
		fmt.Fprintf(ctx, "OK")
	case "/get-op2":
		select_db_rand(h.Dbs, int(times), total)
		fmt.Fprintf(ctx, "OK")
	case "/get-op3":
		select_file("/var/db1_test", int(times), total)
		fmt.Fprintf(ctx, "OK")
	case "/get-op4":
		h.select_memory(int(times), total)
		fmt.Fprintf(ctx, "OK")
	case "/put-op1":
		for i := 0; i < 10; i++ {
			escribir_db(db)
		}
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
	*/
	
}
/*
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
func escribir_db(db int64){

	c := 0
	now := time.Now()
	dbs := getsqlite(int(db))
	for i := 0; i < 10000; i++ {
		if add_txt_db(dbs) {
			c++
		}
	}
	elapsed := time.Since(now)
	fmt.Printf("WRITE %v DB %v en [%v] [%s] c/u\n", c, db, elapsed, time_cu(elapsed, c))
	
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
func add_txt_db(db *sql.DB) bool {
	str := []byte("{\"Id\":1,\"Data\":{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}}")
	stmt, err := db.Prepare("INSERT INTO contents(content) values(?)")
	if err != nil {
		fmt.Println(err)
		return false
	}
	res, err := stmt.Exec(string(str))
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
		return false
	}
	return true
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
func printelaped(start time.Time, str string) {
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
*/
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
func divmod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}
func getdbid(num, base int64) (db, id int64) {
	c, n := divmod(num-base, 100000)
	return c, n
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
func random(max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max))
	return n.Int64()
}
func read_int64(data []byte) int64 {
    var x int64
    for _, c := range data {
        x = x * 10 + int64(c - '0')
    }
    return x
}
func printelaped(start time.Time, str string) {
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
	//return time.Now()
}
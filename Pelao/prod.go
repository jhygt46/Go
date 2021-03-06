package main

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fasthttp/router"
	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
)

type Response struct {
	Op     uint8  `json:"Op"`
	Msg    string `json:"Msg"`
	Reload int    `json:"Reload"`
	Page   string `json:"Page"`
	Tipo   string `json:"Tipo"`
	Titulo string `json:"Titulo"`
	Texto  string `json:"Texto"`
}
type Giros struct {
	Titulo string `json:"Titulo"`
}
type Config struct {
	Tiempo time.Duration `json:"Tiempo"`
}
type MyHandler struct {
	Conf Config `json:"Conf"`
}

type TemplateConf struct {
	Titulo          string  `json:"Titulo"`
	SubTitulo       string  `json:"SubTitulo"`
	SubTitulo2      string  `json:"SubTitulo"`
	FormId          int     `json:"FormId"`
	FormAccion      string  `json:"FormAccion"`
	FormNombre      string  `json:"FormNombre"`
	FormDescripcion string  `json:"FormDescripcion"`
	TituloLista     string  `json:"TituloLista"`
	PageMod         string  `json:"PageMod"`
	DelAccion       string  `json:"DelAccion"`
	DelObj          string  `json:"DelObj"`
	Lista           []Lista `json:"FormDescripcion"`
}
type Lista struct {
	Id     int    `json:"Id"`
	Nombre string `json:"Nombre"`
}
type Data struct {
	Nombre string `json:"Nombre"`
}

var (
	imgPrefix  = []byte("/img/")
	imgHandler = fasthttp.FSHandler("/var/Go/Pelao/img", 1)

	cssPrefix  = []byte("/css/")
	cssHandler = fasthttp.FSHandler("/var/Go/Pelao/css", 1)

	jsPrefix  = []byte("/js/")
	jsHandler = fasthttp.FSHandler("/var/Go/Pelao/js", 1)
)

func main() {

	pass := &MyHandler{Conf: Config{}}
	con := context.Background()
	con, cancel := context.WithCancel(con)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()
	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
					pass.Conf.init()
				case os.Interrupt:
					cancel()
					os.Exit(1)
				}
			case <-con.Done():
				log.Printf("Done.")
				os.Exit(1)
			}
		}
	}()
	go func() {
		r := router.New()
		r.GET("/", Index)
		r.GET("/css/{name}", Css)
		r.GET("/js/{name}", Js)
		r.GET("/img/{name}", Img)
		r.GET("/pages/{name}", Pages)
		r.POST("/login", Login)
		r.POST("/save", Save)
		r.POST("/delete", Delete)
		r.POST("/Salir", Salir)
		fasthttp.ListenAndServe(":80", r.Handler)
	}()
	if err := run(con, pass, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func Save(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	id := Read_uint32bytes(ctx.FormValue("id"))
	resp := Response{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	switch string(ctx.FormValue("accion")) {
	case "guardar_empresa":
		nombre := string(ctx.FormValue("nombre"))
		if id == 0 {
			resp = InsertEmpresa(db, nombre)
		}
		if id > 0 {
			resp = UpdateEmpresa(db, nombre, id)
		}
	default:

	}

	json.NewEncoder(ctx).Encode(resp)
}
func Delete(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	id := Read_uint32bytes(ctx.FormValue("id"))
	resp := Response{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	switch string(ctx.FormValue("accion")) {
	case "borrar_empresa":

		resp = BorrarEmpresa(db, id)

	default:

	}

	json.NewEncoder(ctx).Encode(resp)
}
func Login(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	resp := Response{Op: 2}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	user := string(ctx.PostArgs().Peek("user"))

	res, err := db.Query("SELECT id_usr, pass FROM usuarios WHERE user = ?", user)
	defer res.Close()
	ErrorCheck(err)

	if res.Next() {

		var id_usr int
		var pass string
		err := res.Scan(&id_usr, &pass)
		ErrorCheck(err)

		if pass == GetMD5Hash(ctx.PostArgs().Peek("pass")) {

			resp.Op = 1
			resp.Msg = ""
			cookie := randSeq(32)

			stmt, err := db.Prepare("INSERT INTO sesiones(cookie, id_usr) VALUES(?,?)")
			ErrorCheck(err)
			stmt.Exec(cookie, id_usr)

			authcookie := CreateCookie("cu", cookie, 26280)
			ctx.Response.Header.SetCookie(authcookie)

		} else {
			resp.Msg = "Usuario Contrase??a no existen"
		}

	} else {
		resp.Msg = "Usuario Contrase??a no existen"
	}

	json.NewEncoder(ctx).Encode(resp)
}
func Pages(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("text/html; charset=utf-8")
	name := ctx.UserValue("name")
	id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))

	switch name {
	case "crear":

		t, err := TemplatePage(fmt.Sprintf("html/%s.php", name))
		ErrorCheck(err)

		obj := GetTemplateConf("Titulo", "Subtitulo", "Subtitulo2", "Titulo Lista", "guardar_empresa", fmt.Sprintf("/pages/%s", name), "borrar_empresa", "Empresa")
		lista, found := GetEmpresas()
		if found {
			obj.Lista = lista
		}

		if id > 0 {
			aux, found := GetEmpresa(id)
			if found {
				obj.FormNombre = aux.Nombre
				obj.FormId = id
			}
		} else {
			obj.FormId = 0
		}

		err = t.Execute(ctx, obj)
		ErrorCheck(err)

	default:
		ctx.NotFound()
	}
}
func Index(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html; charset=utf-8")
	token := string(ctx.Request.Header.Cookie("cu"))
	if GetUser(token) {
		fmt.Fprintf(ctx, showFile("html/inicio.php"))
	} else {
		fmt.Fprintf(ctx, showFile("html/login.php"))
	}
}
func Salir(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.DelCookie("cu")
}
func Js(ctx *fasthttp.RequestCtx) {
	jsHandler(ctx)
}
func Css(ctx *fasthttp.RequestCtx) {
	cssHandler(ctx)
}
func Img(ctx *fasthttp.RequestCtx) {
	imgHandler(ctx)
}

// FUNCTION DB //
func GetMySQLDB() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/mydatabase")
	return
}
func GetUser(token string) bool {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT t1.id_usr FROM usuarios t1, sesiones t2 WHERE t2.cookie = ? AND t2.id_usr=t1.id_usr", token)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {
		return true
	} else {
		return false
	}
}
func Permisos(token string, n int) bool {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT * FROM usuarios t1, sesiones t2, usuario_perfil t3, perfil_tarea t4 WHERE t2.cookie = ? AND t2.id_usr=t1.id_usr AND t1.id_usr=t3.id_usr AND t3.id_per=t4.id_per AND t4.id_tar=?", token, n)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {
		return true
	} else {

		res2, err2 := db.Query("SELECT * FROM usuarios t1, sesiones t2, usuario_tarea t3 WHERE t2.cookie = ? AND t2.id_usr=t1.id_usr AND t1.id_usr=t3.id_usr AND t3.id_tar=?", token, n)
		defer res2.Close()
		if err2 != nil {
			log.Fatal(err2)
		}
		if res2.Next() {
			return true
		} else {
			return false
		}

	}
}
func GetEmpresa(id int) (Data, bool) {

	data := Data{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT nombre FROM empresa WHERE id_emp = ? AND eliminado = ?", id, cn)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {

		var nombre string
		err := res.Scan(&nombre)
		if err != nil {
			log.Fatal(err)
		}
		data.Nombre = nombre
		return data, true

	} else {
		return data, false
	}
}
func GetEmpresas() ([]Lista, bool) {

	data := []Lista{}
	b := false

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT id_emp, nombre FROM empresa WHERE eliminado = ?", cn)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}

	var id int
	var nombre string

	for res.Next() {

		err := res.Scan(&id, &nombre)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Lista{Id: id, Nombre: nombre})
		b = true

	}
	return data, b
}

func InsertEmpresa(db *sql.DB, nombre string) Response {

	resp := Response{}
	stmt, err := db.Prepare("INSERT INTO empresa (nombre) VALUES (?)")
	ErrorCheck(err)
	stmt.Exec(nombre)
	if err == nil {
		resp.Op = 1
		resp.Reload = 1
		resp.Page = "crear"
		resp.Msg = "Empresa ingresada correctamente"
	} else {
		resp.Op = 2
		resp.Msg = "La Empresa no pudo ser ingresada"
	}
	return resp
}
func UpdateEmpresa(db *sql.DB, nombre string, id int) Response {

	resp := Response{}
	stmt, err := db.Prepare("UPDATE empresa SET nombre = ? WHERE id_emp = ?")
	ErrorCheck(err)
	_, e := stmt.Exec(nombre, id)
	ErrorCheck(e)
	if e == nil {
		resp.Op = 1
		resp.Reload = 1
		resp.Page = "crear"
		resp.Msg = "Empresa actualizada correctamente"
	} else {
		resp.Op = 2
		resp.Msg = "La Empresa no pudo ser actualizada"
	}
	return resp
}
func BorrarEmpresa(db *sql.DB, id int) Response {

	del := 1
	resp := Response{}
	stmt, err := db.Prepare("UPDATE empresa SET eliminado = ? WHERE id_emp = ?")
	ErrorCheck(err)
	_, e := stmt.Exec(del, id)
	ErrorCheck(e)
	if e == nil {
		resp.Tipo = "success"
		resp.Reload = 1
		resp.Page = "crear"
		resp.Titulo = "Empresa eliminada"
		resp.Texto = "Empresa eliminada correctamente"
	} else {
		resp.Tipo = "error"
		resp.Titulo = "Error al eliminar empresa"
		resp.Texto = "La empresa no pudo ser eliminada"
	}
	return resp
}

// FUNCTION DB //

// DAEMON //
func (h *MyHandler) StartDaemon() {
	h.Conf.Tiempo = 2000 * time.Second
}
func (c *Config) init() {
	var tick = flag.Duration("tick", 1*time.Second, "Ticking interval")
	c.Tiempo = *tick
}
func run(con context.Context, c *MyHandler, stdout io.Writer) error {
	c.Conf.init()
	log.SetOutput(os.Stdout)
	for {
		select {
		case <-con.Done():
			return nil
		case <-time.Tick(c.Conf.Tiempo):
			c.StartDaemon()
		}
	}
}

// DAEMON //
func TemplatePage(v string) (*template.Template, error) {

	t, err := template.ParseFiles(v)
	if err != nil {
		log.Print(err)
		return t, err
	}
	return t, nil
}
func Read_uint32bytes(data []byte) int {
	var x int
	for _, c := range data {
		x = x*10 + int(c-'0')
	}
	return x
}
func GetMD5Hash(text []byte) string {
	hasher := md5.New()
	hasher.Write(text)
	return hex.EncodeToString(hasher.Sum(nil))
}
func CreateCookie(key string, value string, expire int) *fasthttp.Cookie {
	if strings.Compare(key, "") == 0 {
		key = "GoLog-Token"
	}
	fmt.Println("CreateCookie | Key: ", key, " | Val: ", value)
	authCookie := fasthttp.Cookie{}
	authCookie.SetKey(key)
	authCookie.SetValue(value)
	authCookie.SetMaxAge(expire)
	authCookie.SetHTTPOnly(true)
	authCookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	return &authCookie
}
func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func showFile(file string) string {

	dat, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(dat)
}
func ErrorCheck(e error) {
	if e != nil {
		fmt.Println("ERROR:", e)
	}
}
func GetTemplateConf(titulo string, subtitulo string, subtitulo2 string, titulolista string, formaccion string, pagemod string, delaccion string, delobj string) TemplateConf {
	return TemplateConf{Titulo: titulo, SubTitulo: subtitulo, SubTitulo2: subtitulo2, TituloLista: titulolista, FormAccion: formaccion, PageMod: pagemod, DelAccion: delaccion, DelObj: delobj}
}

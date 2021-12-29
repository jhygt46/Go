package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/valyala/fasthttp"

	//"resource/lang"
	//"resource/utils"
	"resource/initserver"
	"resource/kubernet"
)

type Config struct {
	Tiempo       time.Duration `json:"Tiempo"`
	UltimoCambio time.Time     `json:"UltimoCambio"`
}
type MyHandler struct {
	Conf     Config             `json:"Conf"`
	Kubernet *kubernet.Kubernet `json:"Kubernet"`
}

func CreateDb(files []kubernet.Archivo) {

	total := 1000000

	for _, v := range files {
		db, err := getsqlite(v.File)
		if err == nil {
			add_db(db, total)
		}
	}

}

func main() {

	kub := kubernet.Kubernet{}
	kub.Servers = make(map[string]*kubernet.Server, 0)
	kub.Archivos = []kubernet.Archivo{
		kubernet.Archivo{
			Tipo:   1,
			File:   "filtrodb0",
			Ip:     "18.118.187.180",
			Rango1: 1,
			Rango2: 1000000,
		},
		kubernet.Archivo{
			Tipo:   1,
			File:   "filtrodb1",
			Ip:     "18.118.187.180",
			Rango1: 1000001,
			Rango2: 2000000,
		},
	}

	CreateDb(kub.Archivos)

	kub.Configuracion = kubernet.Configuracion{Ip: "18.118.187.180", Port: "8600", UltimoCambio: time.Now()}
	filtros := kubernet.Servicio{
		Tipo:   1,
		Nombre: "filtro",
		Valor:  "/filtro/",
		ListadeBackends: []kubernet.ListadeBackend{
			kubernet.ListadeBackend{
				Activo: true,
				Backends: []kubernet.Backend{
					kubernet.Backend{
						Acls: []kubernet.Acl{
							kubernet.Acl{Param: "id", Tipo: 1, Valor1: 1000, Valor2: 1500000},
						},
						Servers: []kubernet.ServerId{},
					},
					kubernet.Backend{
						Acls: []kubernet.Acl{
							kubernet.Acl{Param: "id", Tipo: 1, Valor1: 1, Valor2: 100},
						},
						Servers: []kubernet.ServerId{},
					},
				},
			},
			kubernet.ListadeBackend{
				Activo: false,
				Backends: []kubernet.Backend{
					kubernet.Backend{
						Acls: []kubernet.Acl{
							kubernet.Acl{Param: "id", Tipo: 1, Valor1: 1, Valor2: 100},
						},
						Servers: []kubernet.ServerId{},
					},
					kubernet.Backend{
						Acls: []kubernet.Acl{
							kubernet.Acl{Param: "id", Tipo: 1, Valor1: 1, Valor2: 100},
						},
						Servers: []kubernet.ServerId{},
					},
				},
			},
		},
	}

	kub.Servicios = []kubernet.Servicio{filtros}
	pass := &MyHandler{Conf: Config{UltimoCambio: time.Now()}, Kubernet: &kub}

	pass.AddServer(0, 0, 0)

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
		fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)
	}()
	if err := run(con, pass, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) == "POST" {
		params := ctx.PostBody()
		switch string(ctx.Path()) {
		case "/newconf":
			res := initserver.ReqInitServer{}
			if err := json.Unmarshal(params, &res); err == nil {
				json.NewEncoder(ctx).Encode(h.Conf.UltimoCambio)
			} else {
				fmt.Println(err)
			}
		case "/init":
			req := initserver.ReqInitServer{}
			if err := json.Unmarshal(params, &req); err == nil {
				json.NewEncoder(ctx).Encode(h.InitServer(req))
			} else {
				fmt.Fprintf(ctx, "")
				fmt.Println("ERROR INIT")
				fmt.Println(err)
			}
		case "/status":
			req := initserver.ReqInitServer{}
			if err := json.Unmarshal(params, &req); err == nil {
				res, err := initserver.Status(fmt.Sprintf("http://%s/status", req.Ip), initserver.ReqStatus{Token: req.Token})
				if err == nil {
					h.InitStatus(res, req.Id)
					fmt.Fprintf(ctx, "")
				} else {
					fmt.Fprintf(ctx, "")
					fmt.Println("ERROR STATUS")
					fmt.Println(err)
				}
			} else {
				fmt.Fprintf(ctx, "")
				fmt.Println("ERROR UNMARSHAL STATUS")
				fmt.Println(err)
			}
		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}
}
func (h *MyHandler) InitServer(req initserver.ReqInitServer) initserver.ResInitServer {

	res := initserver.ResInitServer{}

	if server, found := h.Kubernet.Servers[req.Id]; found {

		res.Encontrado = true
		res.Consulname = fmt.Sprintf("cn%s%d%d", h.Kubernet.Servicios[server.PosicionServicio].Nombre, server.PosicionListaBackend, server.PosicionBackend)
		res.Consulhost = h.Kubernet.Configuracion.Ip + ":" + h.Kubernet.Configuracion.Port

		for _, arch := range h.Kubernet.Archivos {
			if arch.Tipo == h.Kubernet.Servicios[server.PosicionServicio].Tipo {
				for _, acl := range h.Kubernet.Servicios[server.PosicionServicio].ListadeBackends[server.PosicionListaBackend].Backends[server.PosicionBackend].Acls {
					if acl.Valor1 >= arch.Rango1 || acl.Valor2 <= arch.Rango2 {
						if fileexist(res.Files, arch) {
							res.Files = append(res.Files, initserver.File{File: arch.File, Ip: arch.Ip})
						}
					}
				}
			}
		}

		//res.AutoCache = true;
		//res.ListaCache = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
		res.TotalCache = 300000

	} else {
		res.Encontrado = false
	}

	return res
}
func (h *MyHandler) InitStatus(req initserver.ResStatus, Id string) {

	if server, found := h.Kubernet.Servers[Id]; found {

		if !server.Iniciado.Init {
			h.Kubernet.Servers[Id].Ip = "18.92.167.230"
			h.Kubernet.Servers[Id].Iniciado.Consul = req.Consul
			h.Kubernet.Servers[Id].Iniciado.Scp = req.Scp
			h.Kubernet.Servers[Id].Iniciado.Init = req.Init
		}
		if len(req.Cpu) > 0 {
			h.Kubernet.Servers[Id].Cpu = initserver.GetCpuPonderacion(req.Cpu)
		}
		if len(req.Memory) > 0 {
			h.Kubernet.Servers[Id].Memory = initserver.GetMemoryPonderacion(req.Memory)
		}
		h.Kubernet.Servers[Id].DiskMb = req.SizeMb
		fmt.Println(h.Kubernet.Servers[Id])

	}
}
func (h *MyHandler) AddServer(pos_serv int, pos_lista int, pos_bckn int) {
	Id := "i-0c3d0fdd5f1459610"
	h.Kubernet.Servicios[pos_serv].ListadeBackends[pos_lista].Backends[pos_bckn].Servers = append(h.Kubernet.Servicios[pos_serv].ListadeBackends[pos_lista].Backends[pos_bckn].Servers, kubernet.ServerId{Id: Id})
	h.Kubernet.Servers[Id] = &kubernet.Server{
		PosicionServicio:     pos_serv,
		PosicionListaBackend: pos_lista,
		PosicionBackend:      pos_bckn,
	}
}
func (h *MyHandler) DelServer(Id string) {
	if server, found := h.Kubernet.Servers[Id]; found {
		delete(h.Kubernet.Servers, Id)
		for i, v := range h.Kubernet.Servicios[server.PosicionServicio].ListadeBackends[server.PosicionListaBackend].Backends[server.PosicionBackend].Servers {
			if v.Id == Id {
				h.Kubernet.Servicios[server.PosicionServicio].ListadeBackends[server.PosicionListaBackend].Backends[server.PosicionBackend].Servers = kubernet.RemoveServerId(h.Kubernet.Servicios[server.PosicionServicio].ListadeBackends[server.PosicionListaBackend].Backends[server.PosicionBackend].Servers, i)
			}
		}
	} else {
		fmt.Println("DEL NOT FOUND")
	}
}

// DAEMON //
func (h *MyHandler) StartDaemon() {

	h.Conf.Tiempo = 20 * time.Second
	fmt.Println("DAEMON")
	/*
		hpconf := haproxy.Create_config_file(*h.Kubernet)
		err := os.WriteFile("haproxy.cfg", []byte(hpconf), 0644)
		if err == nil {
			fmt.Println("CREATE HAPROXY.CFG")
		}
	*/

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

func (h *MyHandler) AddBlockIp(time time.Time, ip string) {
	h.Kubernet.IpBlocks = append(h.Kubernet.IpBlocks, kubernet.IpBlocks{Tiempo: time, Ip: ip})
}
func fileexist(files []initserver.File, file kubernet.Archivo) bool {
	for _, v := range files {
		if v.File == file.File {
			return false
		}
	}
	return true
}

func add_db(db *sql.DB, total int) {

	str1 := []byte("{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}")
	str := string(str1)
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.
	stmt, err := tx.Prepare("INSERT INTO filtros (filtro, cache) VALUES(?, ?)")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.
	for i := 0; i < total; i++ {
		if _, err := stmt.Exec(str, i); err != nil {
			fmt.Println(err)
		}
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}

}
func getsqlite(dbn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "/var/db/"+dbn)
	if err == nil {
		stmt, err := db.Prepare(`create table if not exists filtros (id integer not null primary key autoincrement,filtro text, cache integer)`)
		if err != nil {
			fmt.Println("err1")
			fmt.Println(err)
			return db, err
		}
		stmt.Exec()
		return db, nil
	} else {
		fmt.Println("err2")
		fmt.Println(err)
		return db, err
	}
}

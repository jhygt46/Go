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

	//"resource/consul"
	//"resource/lang"
	"resource/initserver"
	"resource/scp"
	"resource/utils"
)

type Count struct {
	Cache          int32     `json:"Cache"`
	Db             int32     `json:"Db"`
	UltimaMedicion time.Time `json:"UltimaMedicion"`
	TotalBytes     uint64    `json:"TotalBytes"`
}
type Daemon struct {
	Tiempo       time.Duration `json:"Tiempo"`
	TiempoMemory time.Time     `json:"TiempoMemory"`
	TiempoDisk   time.Time     `json:"TiempoDisk"`
	TiempoCpu    time.Time     `json:"TiempoCpu"`
}
type InfoServer struct {
	Id             string `json:"Id"`
	Ip             string `json:"Ip"`
	Token          string `json:"Token"`
	CacheCapicidad int32  `json:"CacheCapicidad"`
	CacheCount     int32  `json:"CacheCount"`
	StopCache      bool   `json:"StopCache"`
}
type MyHandler struct {
	StatusServer initserver.ResStatus `json:"StatusServer"`
	Count        Count                `json:"Count"`
	Daemon       Daemon               `json:"Daemon"`
	InfoServer   InfoServer           `json:"Info"`
	Cache        map[uint32]Filtro    `json:"Cache"`
}

type Filtro struct {
	C []Campos `json:"C"`
	E []Evals  `json:"E"`
}
type Campos struct {
	T int      `json:"T"`
	N string   `json:"N"`
	V []string `json:"V"`
}
type Evals struct {
	T int    `json:"T"`
	N string `json:"N"`
}

func main() {

	Id := utils.GetInstanceMeta("instance-id")
	Ip := initserver.LocalIP()

	fmt.Printf("Id:%s / Ip:%s\n", Id, Ip)

	pass := &MyHandler{
		Cache:        make(map[uint32]Filtro),
		Daemon:       Daemon{TiempoMemory: time.Now(), TiempoDisk: time.Now(), TiempoCpu: time.Now()},
		Count:        Count{Cache: 0, Db: 0, UltimaMedicion: time.Now(), TotalBytes: 0},
		StatusServer: initserver.ResStatus{SizeMb: 0, Memory: make([]initserver.StatusMemory, 0), Cpu: make([]initserver.StatusCpu, 0), Consul: false, Scp: false, Init: false},
		InfoServer:   InfoServer{Id: Id, Ip: Ip, Token: "", CacheCount: 0, StopCache: false},
	}

	init, err := initserver.Init("http://172.31.32.244/init", initserver.ReqInitServer{Id: pass.InfoServer.Id, Ip: pass.InfoServer.Ip})
	if err == nil {

		if init.Encontrado {

			fmt.Printf("SERVIDOR ENCONTRADO\n")
			pass.StatusServer.Init = true
			pass.InfoServer.CacheCapicidad = init.TotalCache

			pass.StatusServer.Scp = true
			for _, v := range init.Files {
				err := scp.CopyFile(v.Ip, "/var/db/"+v.File, "/var/db/"+v.File)
				if err != nil && pass.StatusServer.Scp {
					pass.StatusServer.Scp = false
				}
				if err != nil {
					fmt.Println(err)
				}
				if err == nil {
					pass.AddCache(v.File)
				}
			}
			/*
				if consul.ConsulRegisters(init.Consulname, init.Consulhost) {
					pass.StatusServer.Consul = true
					fmt.Printf("CONSUL REGISTER\n")
				} else {
					fmt.Printf("ERROR CONSUL\n")
				}
			*/
		} else {
			fmt.Printf("SERVIDOR NO ENCONTRADO\n")
		}

	} else {
		fmt.Printf("ERROR INIT REQUEST\n")
		fmt.Println(err)
	}

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
					pass.Daemon.init()
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
	if string(ctx.Method()) == "GET" {
		switch string(ctx.Path()) {
		case "/filtro":

			ctx.Response.Header.Set("Content-Type", "application/json")
			id := utils.Read_uint32(ctx.QueryArgs().Peek("id"))

			if res, found := h.Cache[id]; found {
				//fmt.Fprintf(ctx, res)
				json.NewEncoder(ctx).Encode(res)
				h.Count.Cache++
			} else {
				h.Count.Db++

				// BUSCAR EN BASE DE DATOS
				if h.InfoServer.CacheCapicidad > h.InfoServer.CacheCount {
					// GUARDAR EN CACHE
					h.InfoServer.CacheCount++
				}
				// MOSTRAR DATA

			}

		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}
	if string(ctx.Method()) == "POST" {
		switch string(ctx.Path()) {
		case "/status":
			params := ctx.PostBody()
			req := initserver.ReqStatus{}
			if err := json.Unmarshal(params, &req); err == nil {
				if req.Token == h.InfoServer.Token {
					h.InfoServer.Token = ""
					json.NewEncoder(ctx).Encode(h.StatusServer)
				}
			} else {
				fmt.Println(err)
				ctx.Error("Not Found", fasthttp.StatusNotFound)
			}
		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}
}

// DAEMON //
func (h *MyHandler) StartDaemon() {

	fmt.Println("TOTAL BYTES:", h.Count.TotalBytes/1024/1024)
	send := false
	h.Daemon.Tiempo = 5 * time.Second

	if h.InfoServer.Token == "" {
		h.InfoServer.Token = initserver.RandStringBytes(10)
	}

	total := int32(time.Since(h.Count.UltimaMedicion) / time.Millisecond)
	totalcache := h.Count.Cache / total
	totaldb := h.Count.Db / total
	h.ResetCount()

	if totalcache+totaldb*2 > 14 || time.Now().After(h.Daemon.TiempoCpu) {

		//statuscpu := initserver.StatusCpu{CountCacheperMilli: totalcache, CountDbperMilli: totaldb, Fecha: time.Now(), CpuUsage: 10, IdleTicks: 10, TotalTicks: 10}
		statuscpu := initserver.GetMonitoringsCpu(totalcache, totaldb)
		u, err := json.Marshal(statuscpu)
		if err == nil {
			fmt.Println(string(u))
		}

		if statuscpu.CpuUsage > 70 {
			send = true
			if len(h.StatusServer.Cpu) > 9 {
				h.StatusServer.Cpu = initserver.RemoveIndexCpu(h.StatusServer.Cpu, 0)
			}
			h.StatusServer.Cpu = append(h.StatusServer.Cpu, statuscpu)
		}
		h.Daemon.TiempoCpu = h.Daemon.TiempoCpu.Add(30 * time.Second)
	}

	if !h.InfoServer.StopCache || time.Now().After(h.Daemon.TiempoMemory) {

		//statusmemory := initserver.StatusMemory{Fecha: time.Now(), Alloc: 10, TotalAlloc: 10, Sys: 10, NumGC: 10}
		statusmemory := initserver.PrintMemUsage()
		if statusmemory.Alloc > 300 {
			send = true
			if len(h.StatusServer.Memory) > 9 {
				h.StatusServer.Memory = initserver.RemoveIndexMem(h.StatusServer.Memory, 0)
			}
			h.StatusServer.Memory = append(h.StatusServer.Memory, statusmemory)
			h.InfoServer.StopCache = true
		} else {
			if h.InfoServer.CacheCount == h.InfoServer.CacheCapicidad {
				h.InfoServer.CacheCapicidad += 1000
				h.InfoServer.StopCache = false
			}
			h.Daemon.TiempoMemory = h.Daemon.TiempoMemory.Add(30 * time.Second)
		}

	}

	if time.Now().After(h.Daemon.TiempoDisk) {
		size, err := initserver.DirSize("C:/Allin/GoFinal/Filtros")
		if err == nil {
			if size > 3000 {
				send = true
				h.StatusServer.SizeMb = size
			}
		}
		h.Daemon.TiempoDisk = h.Daemon.TiempoDisk.Add(300 * time.Second)
	}
	if send {
		fmt.Println("ENVIANDO STATUS")
		_, err := initserver.Status("http://localhost:81/status", initserver.ReqStatus{Id: h.InfoServer.Id, Ip: h.InfoServer.Ip, Token: h.InfoServer.Token})
		if err != nil {
			fmt.Println(err)
		}
	}
}
func (c *Daemon) init() {
	var tick = flag.Duration("tick", 1*time.Second, "Ticking interval")
	c.Tiempo = *tick
}
func run(con context.Context, c *MyHandler, stdout io.Writer) error {
	c.Daemon.init()
	log.SetOutput(os.Stdout)
	for {
		select {
		case <-con.Done():
			return nil
		case <-time.Tick(c.Daemon.Tiempo):
			c.StartDaemon()
		}
	}
}

// DAEMON //

func (h *MyHandler) ResetCount() {
	h.Count.UltimaMedicion = time.Now()
	h.Count.Cache = 0
	h.Count.Db = 0
}
func (h *MyHandler) AddCache(file string) {

	fmt.Print("ADD CACHE:", file)

	db, err := sql.Open("sqlite3", "/var/db/"+file)

	if err == nil {
		rows, err := db.Query("SELECT id, filtro FROM filtros LIMIT 100000")
		if err == nil {
			defer rows.Close()
			var id uint32
			var filtro string
			for rows.Next() {
				err := rows.Scan(&id, &filtro)
				if err == nil {

					data := Filtro{}
					if err := json.Unmarshal([]byte(filtro), &data); err == nil {
						h.Cache[id] = data
					}

					//h.Cache[id] = filtro
					//h.Count.TotalBytes = h.Count.TotalBytes + uint64(len(filtro))
				} else {
					fmt.Print("ERR SCAN:")
					fmt.Println(err)
				}
			}
		} else {
			fmt.Print("ERR SELECT TABLE FILTROS:")
			fmt.Println(err)
		}
	} else {
		fmt.Print("ERR CONNECT DB:", file)
		fmt.Println(err)
	}
}

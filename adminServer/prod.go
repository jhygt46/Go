package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/namsral/flag"
	"encoding/json"
	"io/ioutil"
)

const defaultTick = 60 * time.Second

type Daemon struct {
	Servicios []Servicio `json:"Servicios"`
}
type Servicio struct {
	Servers []Server `json:"Servers"`
	Alertas []Alerta `json:"Alertas"`
	Info Infoservicio `json:"Info"`
}
type Infoservicio struct {
	Nombre string `json:"Nombre"`
	Fecha time.Time `json:"Fecha"`
}
type Alerta struct {
	Tipo int `json:"Tipo"`
	Fecha time.Time `json:"Fecha"`
}
type Server struct {
	Ip string `json:"Ip"`
	Nombre string `json:"Nombre"`
	Acciones []AccionServer `json:"Acciones"`
}
type AccionServer struct {
	Fecha time.Time `json:"Fecha"`
	Nombre string `json:"Nombre"`
	Tipo int `json:"Tipo"`
}
type config struct {
	contentType string
	server      string
	statusCode  int
	tick        time.Duration
	url         string
	userAgent   string
}
func run(ctx context.Context, c *config, stdout io.Writer, daemon *Daemon) error {
	c.init(os.Args)
	log.SetOutput(os.Stdout)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(c.tick):
			daemon.StartDaemon()
		}
	}
}
func (c *config) init(args []string) error {

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.String(flag.DefaultConfigFlagname, "", "Path to config file")

	var (
		statusCode  = flags.Int("status", 200, "Response HTTP status code")
		tick        = flags.Duration("tick", defaultTick, "Ticking interval")
		server      = flags.String("server", "", "Server HTTP header value")
		contentType = flags.String("content_type", "", "Content-Type HTTP header value")
		userAgent   = flags.String("user_agent", "", "User-Agent HTTP header value")
		url         = flags.String("url", "", "Request URL")
	)

	p := []string{"-config=config.conf"}

	if err := flags.Parse(p); err != nil {
		return err
	}

	c.statusCode = *statusCode
	c.tick = *tick
	c.server = *server
	c.contentType = *contentType
	c.userAgent = *userAgent
	c.url = *url

	return nil
}
func main() {

	dae := readFile("daemon.json")
	Daemon := &dae

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
	c := &config{}

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
					c.init(os.Args)
				case os.Interrupt:
					cancel()
					os.Exit(1)
				}
			case <-ctx.Done():
				log.Printf("Done.")
				os.Exit(1)
			}
		}
	}()
	
	if err := run(ctx, c, os.Stdout, Daemon); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}

func (d *Daemon) StartDaemon() {
	
	//s := Servicio{ Servers: []Server{Server{ Ip: "140.0.0.1", Nombre: "Server1", Acciones: []AccionServer{AccionServer{ Nombre: "NELSON", Tipo: 1 }} }}, Info: Infoservicio{ Nombre: "BUE" }, Alertas: []Alerta{Alerta{ Tipo: 1, Fecha: time.Now() }}}
	//server := Server{ Ip: "140.0.0.1", Nombre: "Server1", Acciones: []AccionServer{AccionServer{ Nombre: "NELSON", Tipo: 1 }} }
	//Daemon.AddServicio(Servicio{ Servers: []Server{server, server}, Info: Infoservicio{ Nombre: "BUE" }, Alertas: []Alerta{Alerta{ Tipo: 1, Fecha: time.Now() }}})
	
	for i := range d.Servicios {

		//**d.AddServer(i, Server{ Ip: "140.0.0.1", Nombre: "Server1"})
		//**d.AddAlerta(i, Alerta{ Tipo: 1 })
		//**d.AddSrvAccion(i, j, AccionServer{ Nombre: "ACC1", Tipo: 1 })
		//**d.RemoveServer(i, 0)
		//info := d.Servicios[i].Info

		/*
		for j := range d.Servicios[i].Servers {
			// LISTA DE SERVERS
			// LISTA ACCIONES SERVER
			for x := range d.Servicios[i].Servers[j].Acciones {
				
				switch d.Servicios[i].Servers[j].Acciones[x].Tipo {
				case 1:
					// http.Get("URL CHECK")
				case 2:
					// http.Get("CHECK PROCESADOR")
				case 3:
					// http.Get("CHECK MEMORIA")
				case 4:
					// http.Get("CHECK OTROS")
				}
			}
		}
		*/
		
		for j := range d.Servicios[i].Alertas {
			// ALERTAS DEL SERVICIO
			switch d.Servicios[i].Alertas[j].Tipo {
			case 1:
				d.Alerta1(i, j)
			case 2:
				d.Alerta2(i, j)
			case 3:
				d.Alerta3(i, j)
			case 4:
				d.Alerta4(i, j)
			}
		}
		//if duration.Milliseconds() < 0 {}
	}
	//saveFile(d)
	//d.RemoveAlerta(0)
	//d.RemoveServer(0)

}
func saveFile(data *Daemon){
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile("daemon.json", file, 0644)
}
func (d *Daemon) Alerta1(servicio int, pos int){
	fmt.Println(d.Servicios[servicio].Alertas[pos])
}
func (d *Daemon) Alerta2(servicio int, pos int){
	fmt.Println(d.Servicios[servicio].Alertas[pos])
}
func (d *Daemon) Alerta3(servicio int, pos int){
	fmt.Println(d.Servicios[servicio].Alertas[pos])
}
func (d *Daemon) Alerta4(servicio int, pos int){
	fmt.Println(d.Servicios[servicio].Alertas[pos])
}
func (d *Daemon) RemoveServer(servicio int, ser int) bool {
	if len(d.Servicios[servicio].Servers) > 0 {
		(*d).Servicios[servicio].Servers = append(d.Servicios[servicio].Servers[:ser], d.Servicios[servicio].Servers[ser+1:]...)
		return true
	}
	return false
}
func (d *Daemon) RemoveAlerta(servicio int, pos int) bool {
	if len(d.Servicios[servicio].Alertas) > 0 {
		(*d).Servicios[servicio].Alertas = append(d.Servicios[servicio].Alertas[:pos], d.Servicios[servicio].Alertas[pos+1:]...)
		return true
	}
	return false
}
func (d *Daemon) AddServicio(item Servicio) bool {
	(*d).Servicios = append(d.Servicios, item)
	return true
}
func (d *Daemon) AddServer(servicio int, item Server) bool {
	(*d).Servicios[servicio].Servers = append(d.Servicios[servicio].Servers, item)
	return true
}
func (d *Daemon) AddSrvAccion(servicio int, server int, item AccionServer) bool {
	(*d).Servicios[servicio].Servers[server].Acciones = append(d.Servicios[servicio].Servers[server].Acciones, item)
	return true
}
func (d *Daemon) AddAlerta(servicio int, item Alerta) bool {
	(*d).Servicios[servicio].Alertas = append(d.Servicios[servicio].Alertas, item)
	return true
}
func readFile(file string) Daemon {

	jsonFile, err := os.Open(file)
	var dae Daemon
    if err != nil {
        return dae
    }
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal(byteValue, &dae)
	return dae

}

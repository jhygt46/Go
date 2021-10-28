package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"flag"
	"time"
	"syscall"
	"context"
	"os/signal"
	"io/ioutil"
	"encoding/json"
    "github.com/valyala/fasthttp"
)

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

type Config struct {
	Id int8 `json:"Id"`
	Fecha time.Time `json:"Fecha"`
	Tiempo time.Duration `json:"Time"`
}
type Daemon struct {
	Servicios []Servicio `json:"Servicios"`
}
type MyHandler struct {
	Conf Config `json:"Conf"`
	Dae *Daemon `json:"Dae"`
}

func main() {

	dae := readFile("daemon.json")
	pass := &MyHandler{ Conf: Config{}, Dae: dae }

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
	fmt.Println(h.Conf)
	fmt.Println(*h.Dae)
	fmt.Fprintf(ctx, "ERROR DDos");
}

func (h *MyHandler) StartDaemon() {

	fmt.Println("DAEMON: ", h.Conf.Tiempo)
	h.Conf.Tiempo = 20 * time.Second

}
func (c *Config) init() {

	var tick = flag.Duration("tick", 5 * time.Second, "Ticking interval")
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

func printelaped(start time.Time, str string){
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
}

func readFile(file string) *Daemon {

	jsonFile, err := os.Open(file)
	var dae Daemon
    if err != nil {
        return &dae
    }
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal(byteValue, &dae)
	return &dae

}
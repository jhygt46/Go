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
    "github.com/valyala/fasthttp"
	
	//LOCAL MODULES//
	//"consul"
	//"lang"
	//"utils"
	//"monitoring"
	//"initserver"
	//"scp"
)
type Config struct {
	Tiempo time.Duration `json:"Tiempo"`
}
type MyHandler struct {
	Conf Config `json:"Conf"`
}
func main() {

	pass := &MyHandler{ Conf: Config{} }
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
	if string(ctx.Method()) == "GET" {
		switch string(ctx.Path()) {
		case "/":
		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}
}

// DAEMON //
func (h *MyHandler) StartDaemon() {
	
	h.Conf.Tiempo = 2 * time.Second
	fmt.Println("DAEMON")
	//now := time.Now()
	//fmt.Printf("WRITES FILES %v [%s] c/u\n", 1, utils.Time_cu(time.Since(now), 1))
	//fmt.Println(monitoring.GetMonitoringsCpu())
	//fmt.Println(monitoring.PrintMemUsage())
}
func (c *Config) init() {
	var tick = flag.Duration("tick", 1 * time.Second, "Ticking interval")
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
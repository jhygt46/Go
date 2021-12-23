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
	
	//LOCAL MODULES//
	//"consul"
	//"lang"
	//"monitoring"
	//"initserver"
	//"scp"
	//"utils"
	"kubernet"
)
type Config struct {
	Tiempo time.Duration `json:"Tiempo"`
	UltimoCambio time.Time `json:"UltimoCambio"`
}
type MyHandler struct {
	Conf Config `json:"Conf"`
}
func main() {

	//Id := utils.GetInstanceMeta("instance-id")

	res, err := kubernet.NewConfig("http://localhost:81/init", kubernet.ReqInitServer{ Id: "id-6315546788", Ip: "3.14.127.34" })
	if err == nil {
		fmt.Println("res")
		fmt.Println(res)
	}else{
		fmt.Println("err")
		fmt.Println(err)
	}


	pass := &MyHandler{ Conf: Config{ UltimoCambio: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) } }
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
	if err := run(con, pass, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}

// DAEMON //
func (h *MyHandler) StartDaemon() {
	
	

	/*
	res, err := haproxy.NewConfig("http://localhost:81/status", haproxy.PostRequest{ Time: h.Conf.UltimoCambio })
	if err == nil {
		fmt.Println("res")
		fmt.Println(res)
	}else{
		fmt.Println("err")
		fmt.Println(err)
	}
	*/
	h.Conf.Tiempo = 20 * time.Second
	fmt.Println("DAEMON")

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
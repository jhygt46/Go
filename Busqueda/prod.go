package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)
	}()
	if err := run(con, pass, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	var arrcuad []uint16

	if string(ctx.Method()) == "GET" {
		switch string(ctx.Path()) {
		case "/":
			if err := json.Unmarshal(ctx.QueryArgs().Peek("cuads"), &arrcuad); err == nil {
				for _, element := range arrcuad {
					BadgerKey := append(Read_uint32bytes(ctx.QueryArgs().Peek("cat")), int16tobytes(element)...)
					fmt.Println(BadgerKey)
				}
			} else {
				fmt.Println(err)
			}
		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}
}

// DAEMON //
func (h *MyHandler) StartDaemon() {

	h.Conf.Tiempo = 200 * time.Second
	fmt.Println("DAEMON")
	//now := time.Now()
	//fmt.Printf("WRITES FILES %v [%s] c/u\n", 1, utils.Time_cu(time.Since(now), 1))
	//fmt.Println(monitoring.GetMonitoringsCpu())
	//fmt.Println(monitoring.PrintMemUsage())
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

func Read_uint32bytes(data []byte) []byte {
	var x uint32
	for _, c := range data {
		x = x*10 + uint32(c-'0')
	}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, x)
	return b
}

func int16tobytes(i uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, i)
	return b
}

/*
package main

import (
	"encoding/json"
	"fmt"
)

type Empresa1 struct {
	P Info       `json:"P"`
	I []ProdsId  `json:"I"`
	S []ProdsIds `json:"S"`
}
type Info struct {
	A float64 `json:"A"`
	N float64 `json:"N"`
	C float64 `json:"C"`
}
type ProdsId struct {
	I uint64 `json:"I"`
	P uint64 `json:"P"`
}
type ProdsIds struct {
	I uint64     `json:"I"`
	P uint64     `json:"P"`
	C float64    `json:"C"`
	F [][]uint32 `json:"F"`
	N string     `json:"N"`
}

func main() {

	var f [][]uint32
	f = append(f, []uint32{1})
	f = append(f, []uint32{15, 21, 33, 40, 52})
	f = append(f, []uint32{65535})

	emp := Empresa1{
		P: Info{A: 33.0, N: 67, C: 6.75},
		S: []ProdsIds{
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.5, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.7, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
			ProdsIds{I: 34234, N: "BUENA NELSON", P: 75467, C: 6.8, F: f},
		},
		I: []ProdsId{
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
			ProdsId{I: 34234, P: 75467},
		},
	}

	u, _ := json.Marshal(emp)

	fmt.Println(len(u), "BYTES")
	//fmt.Println(len(emp.I))

	for _, v := range emp.S {
		//fmt.Println(k)
		fmt.Println(v.F)
	}
}

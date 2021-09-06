package main

import (
	"os"
	"io"
	"log"
	"fmt"
	"net"
	"flag"
	"math"
	"time"
	//"unsafe"
	"syscall"
	"context"
	"strconv"
	"os/signal"
	"io/ioutil"
	"encoding/json"
	//"github.com/DmitriyVTitov/size"
    "github.com/valyala/fasthttp"
    //"github.com/dgraph-io/ristretto"
)

type ConfigIp struct {
	Ddos bool `json:"Dos"`
	Fecha time.Time `json:"FechaDos"`
	Ipddos []Ipinfo
}
type Ipinfo struct {
	Ip net.IP `json:"Ip"`
	Fecha time.Time `json:"Fecha"`
	Valor int8 `json:"Valor"`
}
type ConfigAuto struct {
	Auto bool `json:"Auto"`
	Fecha time.Time `json:"FechaDos"`
	Lista []Autoinfo
}
type Autoinfo struct {
	Id int64 `json:"Id"`
	Fecha time.Time `json:"Fecha"`
	Valor int8 `json:"Valor"`
}
type Config struct {
	Id int8 `json:"Id"`
	Fecha time.Time `json:"Fecha"`
	Time time.Duration `json:"Time"`
	Count uint32 `json:"Count"`
	MaxCount uint32 `json:"MaxCount"`
	StartCount bool `json:"StartCount"`
	StartCountCache bool `json:"StartCountCache"`
	CountCache uint64 `json:"CountCache"`
	CountFiles uint64 `json:"CountFiles"`
}
type Data struct {
	C uint16 `json:"C"`
	F int64 `json:"F"`
	E int64 `json:"E"`
}
type MyHandler struct {
	minicache map[uint32]*Data
	ConfIp *ConfigIp
	ConfAuto *ConfigAuto
	Conf *Config
}

func main() {

	ipflag := flag.Int("ip", 30, "")
	cantmemflag := flag.Int("cantmem", 3, "")
	cacheautoflag := flag.Int("cacheauto", 60, "")
	flag.Parse()

	pass := &MyHandler {
		minicache: make(map[uint32]*Data, *cantmemflag), 
		Conf: &Config{ Id: 8, Fecha: time.Now(), Count: 0, MaxCount: uint32(*cantmemflag) }, 
		ConfAuto: &ConfigAuto{ Auto: true, Fecha: time.Now(), Lista: make([]Autoinfo, *cacheautoflag) }, 
		ConfIp: &ConfigIp{ Ddos: true, Fecha: time.Now(), Ipddos: make([]Ipinfo, *ipflag) },
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
					pass.Conf.init(os.Args)
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
		fasthttp.ListenAndServe(":81", pass.HandleFastHTTP)
	}()

	if err := run(con, pass, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}
func (h *MyHandler) StartDaemon() {

	fmt.Println(*h.Conf)
	h.Conf.Time = 20 * time.Second

}
func (c *Config) init(args []string) {

	var tick = flag.Duration("tick", 2 * time.Second, "Ticking interval")
	c.Time = *tick

}
func run(con context.Context, c *MyHandler, stdout io.Writer) error {

	c.Conf.init(os.Args)
	log.SetOutput(os.Stdout)
	for {
		select {
		case <-con.Done():
			return nil
		case <-time.Tick(c.Conf.Time):
			c.StartDaemon()
		}
	}

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	fmt.Println(ctx.RemoteIP())

	switch string(ctx.Path()) {
	case "/filtro":

		if !h.ConfIp.Ddos {

			id, err := strconv.Atoi(string(ctx.QueryArgs().Peek("id")))
			if err == nil {
				if res, found := h.minicache[uint32(id)]; found {

					if h.Conf.StartCount {
						h.Conf.CountCache++
					}
					if h.Conf.StartCountCache {
						h.minicache[uint32(id)].C++
					}
					ctx.Response.Header.Set("Content-Type", "application/json")
					json.NewEncoder(ctx).Encode(res)

				}else{
					jsonFile, err := os.Open("../utils/filtros/"+string(ctx.QueryArgs().Peek("id"))+".json")
					if err == nil{

						ctx.Response.Header.Set("Content-Type", "application/json")
						byteValue, _ := ioutil.ReadAll(jsonFile)

						if h.ConfAuto.Auto {
							if h.Conf.Count < h.Conf.MaxCount {
								data := Data{}
								_ = json.Unmarshal(byteValue, &data)
								h.minicache[uint32(id)] = &data
								h.Conf.Count++
							}else{
								h.ConfAuto.Auto = false
								h.Conf.StartCount = true
							}
						}
						if h.Conf.StartCount {
							h.Conf.CountFiles++
						}
						if h.Conf.StartCountCache {
							
						}
						fmt.Fprintf(ctx, string(byteValue))

					}else{
						fmt.Println(err)
						ctx.Error("Not Found", fasthttp.StatusNotFound)
					}
				}
			}else{
				ctx.Error("Not Found", fasthttp.StatusNotFound)
			}

		}else{
			fmt.Fprintf(ctx, "ERROR DDos");
		}

	case "/health":
		fmt.Fprintf(ctx, "OK");
	case "/info":
		json.NewEncoder(ctx).Encode(h.Conf)
	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}
	

}
func (h *MyHandler) ResetCount(){
	for _, element := range h.minicache {
        element.C = 0
    }
}
func (h *MyHandler) registerCache(id int) bool {
	
	for i, v := range h.ConfAuto.Lista {
		if v.Id == int64(id) {
			h.ConfAuto.Lista[i].Valor = newVal(v)
			if h.ConfAuto.Lista[i].Valor > 100 {
				for _, element := range h.minicache {
					if element.C < 100 {
						return true 
					}
				}
			}
		}else{
			fmt.Println(v)
		}
	}
	
	return true

}
func newVal(v Autoinfo) int8 {
	return int8(101)
}
func printelaped(start time.Time, str string){
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
}
func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}
func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%dB", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}
	return fmt.Sprintf(f+"%s", val, suffix)
}
func FileSize(s int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(uint64(s), 1024, sizes)
}
func LocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

/*
fmt.Println(FileSize(int64(size.Of(pass))))
fmt.Println("Memory size of Data", unsafe.Sizeof(Data{}))
fmt.Println("Se crearon: ", *cantmemflag)
*/
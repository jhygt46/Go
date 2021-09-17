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
	"bufio"
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
type Campos struct {
	T int `json:"T"`
	N string `json:"N"`
	V [] string `json:"V"`
}
type Evals struct {
	T int `json:"T"`
	N string `json:"N"`
}
type Data struct {
	C [] Campos `json:"C"`
	E [] Evals `json:"E"`
}
type Filtros struct {
	Id int `json:"Id"`
	Data Data `json:"Data"`
}
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
	Tiempo time.Duration `json:"Time"`
}
type Metricas struct {
	Start bool `json:"Start"`
	Fecha time.Time `json:"Fecha"`
	Count uint64 `json:"Count"`
	CountCache uint64 `json:"CountCache"`
	CountFiles uint64 `json:"CountFiles"`
}
type AutoCache struct {
	Start bool `json:"Start"`
	Count uint64 `json:"Count"`
	TotalCache uint64 `json:"TotalCache"`
}
/*
type Datas struct {
	C uint16 `json:"C"`
	F int64 `json:"F"`
	E int64 `json:"E"`
}
*/
type MyHandler struct {
	minicache map[uint32]*Data
	ConfIp *ConfigIp
	ConfAuto *ConfigAuto
	Conf *Config
	Metricas *Metricas
	AutoCache *AutoCache
}
func main() {

	ipflag := flag.Int("ip", 30, "")
	totalcache := flag.Int("totalcache", 3000, "")
	cacheautoflag := flag.Int("cacheauto", 60, "")
	flag.Parse()

	pass := &MyHandler {
		minicache: make(map[uint32]*Data, *totalcache), 
		Conf: &Config{ Id: 8, Fecha: time.Now() },
		ConfAuto: &ConfigAuto{ Auto: true, Fecha: time.Now(), Lista: make([]Autoinfo, *cacheautoflag) }, 
		ConfIp: &ConfigIp{ Ddos: false, Fecha: time.Now(), Ipddos: make([]Ipinfo, *ipflag) },
		Metricas: &Metricas{  Start: false, Fecha: time.Now(), Count: 0, CountCache: 0, CountFiles: 0 },
		AutoCache: &AutoCache{ Start: true, Count: 0, TotalCache: uint64(*totalcache) },
	}

	file1 := "../utils/cache/cachedata2.json"
	if FileExists(file1) {

		start := time.Now()

		f, err := os.Open(file1)
		if err != nil { log.Fatalf("Error to read [file=%v]: %v", file1, err.Error()) }
		fi, err := f.Stat()
		if err != nil { log.Fatalf("Could not obtain stat, handle error: %v", err.Error()) }

		r := bufio.NewReader(f)
		dec := json.NewDecoder(r)
		i := 0

		dec.Token()
		for dec.More() {
			var m Filtros
			err := dec.Decode(&m)
			if err != nil {
				log.Fatal(err)
			}
			i++
			pass.minicache[uint32(m.Id)] = &m.Data

		}
		dec.Token()

		elapsed := time.Since(start)
		fmt.Printf("CacheData Cantidad [%v] Peso [%s] Tiempo [%v] .\n", i, FileSize(fi.Size()), elapsed)

	}

	file2 := "../utils/cache/cachelist2.json"
	if FileExists(file2) {

		start := time.Now()
		i := 0

		jsonFile, err := os.Open(file2)
		if err == nil{
			var list []int
			byteValue, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &list)
			for _, v := range list {
				if pass.minicache[uint32(v)] == nil {
					jsonFiltro, err2 := os.Open("../utils/filtros/"+strconv.Itoa(v)+".json")
					if err2 == nil {
						byteValueFiltro, _ := ioutil.ReadAll(jsonFiltro)
						data := Data{}
						_ = json.Unmarshal(byteValueFiltro, &data)
						pass.minicache[uint32(v)] = &data
						i++
					}
					defer jsonFiltro.Close()
				}
			}
		}
		defer jsonFile.Close()

		elapsed := time.Since(start)
		fmt.Printf("CacheList Cantidad [%v] Tiempo [%v] .\n", i, elapsed)

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
		fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)
	}()

	if err := run(con, pass, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}
func (h *MyHandler) StartDaemon() {

	fmt.Println("DAEMON")
	h.Conf.Tiempo = 5 * time.Second

}
func (c *Config) init(args []string) {

	var tick = flag.Duration("tick", 5 * time.Second, "Ticking interval")
	c.Tiempo = *tick

}
func run(con context.Context, c *MyHandler, stdout io.Writer) error {

	c.Conf.init(os.Args)
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

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	//fmt.Println(ctx.RemoteIP())

	switch string(ctx.Path()) {
	case "/filtro":

		if !h.ConfIp.Ddos {

			id, err := strconv.Atoi(string(ctx.QueryArgs().Peek("id")))
			if err == nil {
				if res, found := h.minicache[uint32(id)]; found {

					//if h.Metricas.Start { h.Metricas.CountCache++ }
					ctx.Response.Header.Set("Content-Type", "application/json")
					json.NewEncoder(ctx).Encode(res)

				}else{
					jsonFile, err := os.Open("../utils/filtros/"+string(ctx.QueryArgs().Peek("id"))+".json")
					if err == nil{

						ctx.Response.Header.Set("Content-Type", "application/json")
						byteValue, _ := ioutil.ReadAll(jsonFile)
						if h.AutoCache.Start {
							if h.AutoCache.Count < h.AutoCache.TotalCache {
								data := Data{}
								_ = json.Unmarshal(byteValue, &data)
								h.minicache[uint32(id)] = &data
								h.AutoCache.Count++
							}else{
								h.AutoCache.Start = false
							}
						}
						if h.Metricas.Start {
							h.Metricas.CountFiles++
						}
						fmt.Fprintf(ctx, string(byteValue))

					}else{
						fmt.Println(err)
						ctx.Error("Not Found", fasthttp.StatusNotFound)
					}
					defer jsonFile.Close()
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
	case "/metrica":
		if h.Metricas.CountFiles > 0 {
			fmt.Fprintf(ctx, strconv.FormatUint(h.Metricas.CountCache/h.Metricas.CountFiles, 10)+"s")
		}else{
			fmt.Fprintf(ctx, "BUENA")
		}
	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}
	

}
func (h *MyHandler) ResetCount(){
	for _, element := range h.minicache {
		fmt.Println(element)
    }
}
func (h *MyHandler) registerCache(id int) bool {
	
	for i, v := range h.ConfAuto.Lista {
		if v.Id == int64(id) {
			h.ConfAuto.Lista[i].Valor = newVal(v)
			if h.ConfAuto.Lista[i].Valor > 100 {
				for _, element := range h.minicache {
					fmt.Println(element)
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
func FileExists(name string) bool {
    if fi, err := os.Stat(name); err == nil {
        if fi.Mode().IsRegular() {
            return true
        }
    }
    return false
}

/*
fmt.Println(FileSize(int64(size.Of(pass))))
fmt.Println("Memory size of Data", unsafe.Sizeof(Data{}))
fmt.Println("Se crearon: ", *cantmemflag)
*/
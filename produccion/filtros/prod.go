package main

import (
	"os"
	"fmt"
	"net"
	"log"
	"time"
	"strconv"
	"math/big"
	"net/http"
	"io/ioutil"
	"crypto/rand"
	"path/filepath"
	"encoding/json"
    "github.com/valyala/fasthttp"
	"github.com/hashicorp/consul/api"
	"bitbucket.org/bertimus9/systemstat"
)

// MONITORING //
var coresToPegPtr *int64
type stats struct {
	startTime time.Time

	// stats this process
	ProcUptime        float64 //seconds
	ProcMemUsedPct    float64
	ProcCPUAvg        systemstat.ProcCPUAverage
	LastProcCPUSample systemstat.ProcCPUSample `json:"-"`
	CurProcCPUSample  systemstat.ProcCPUSample `json:"-"`

	// stats for whole system
	LastCPUSample systemstat.CPUSample `json:"-"`
	CurCPUSample  systemstat.CPUSample `json:"-"`
	SysCPUAvg     systemstat.CPUAverage
	SysMemK       systemstat.MemSample
	LoadAverage   systemstat.LoadAvgSample
	SysUptime     systemstat.UptimeSample

	// bookkeeping
	procCPUSampled bool
	sysCPUSampled  bool
}
// MONITORING //

type MyHandler struct {
	minicache map[uint32]*Data
}
type Data struct {
	C [] Campos `json:"C"`
	E [] Evals `json:"E"`
}
type Campos struct {
	T int `json:"T"`
	N string `json:"N"`
	V [] string `json:"V"`
}
type Evals struct {
	T int `json:"T"`
	N string `json:"N"`
}

type adminResponse struct {
	Consulname string `json:"Consulname"`
	Consulhost string `json:"Consulip"`
	Cachetipo int8 `json:"Cachetipo"` // 0 AUTOMATICO - 1 LISTA CACHE
	ListaCache []int64 `json:"ListaCache"`
}

type ConsulRegister struct {
	Address                        string
	Name                           string
	Tag                            []string
	Port                           int
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
}

func main() {

	h := &MyHandler{}
	//h.initServer()

	size, err := DirSize("/var/Go")
	if err == nil {
		fmt.Println(size)
	}

	fasthttp.ListenAndServe(":81", h.HandleFastHTTP)
	
}

func (h *MyHandler) initServer() {
	
	res := getUrl("http://localhost/")

	switch res.Cachetipo {
	case 0:

	case 1:
		for _, v := range res.ListaCache {
			fmt.Println("Lista: ", v)
		}
	default:
		
	}

	if consulRegister(res.Consulname, res.Consulhost) {

	}

}
func getUrl(url string) *adminResponse {

	myClient := &http.Client{Timeout: 10 * time.Second}
    r, err := myClient.Get(url)
	var admin adminResponse
    if err != nil {
        return &admin
    }
    defer r.Body.Close()
    json.NewDecoder(r.Body).Decode(&admin)
	return &admin

}
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	switch string(ctx.Path()) {
	case "/filtro":
		ctx.Response.Header.Set("Content-Type", "application/json")
		id := read_int32(ctx.QueryArgs().Peek("id"))
		if res, found := h.minicache[id]; found {
			json.NewEncoder(ctx).Encode(res)
		}else{

		}
	case "/monitoring":
		stats := NewStats()
		stats.GatherStats()
		stats.PrintStats()
	default:
		//ctx.Error("Not Found", fasthttp.StatusNotFound)
		fmt.Fprintf(ctx, "ok");
	}
	
}

// AWS METADATA //
func getInstanceId() string {

	resp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
   	if err != nil {
		log.Fatalln(err)
   	}
	return string(body)
    
}
// AWS METADATA //

// UTILS //
func read_int32(data []byte) uint32 {
    var x uint32
    for _, c := range data {
        x = x * 10 + uint32(c - '0')
    }
    return x
}
// UTILS //

// CONSUL //
func consulRegister(name string, consuladress string) bool {

	s := NewConsulRegister(name, consuladress)
	config := api.DefaultConfig()
	config.Address = s.Address
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Println(err)
	}
	agent := client.Agent()

	IP := LocalIP()
	reg := &api.AgentServiceRegistration{
		 ID: fmt.Sprintf("%v-%v-%v", s.Name, IP, s.Port), // Name of the service node
		 Name: s.Name, // service name
		 Tags: s.Tag, // tag, can be empty
		 Port: s.Port, // service port
		 Address: IP, // Service IP
		 Check: &api.AgentServiceCheck{ // Health Check
			 Interval: s.Interval.String(), // Health check interval
			 HTTP: fmt.Sprintf("http://%s:%d/%s", IP, s.Port, s.Name), // grpc support, address to perform health check, service will be passed to Health.Check function
			 DeregisterCriticalServiceAfter: s.DeregisterCriticalServiceAfter.String(), // Deregistration time, equivalent to expiration time
		},
	}

	if err := agent.ServiceRegister(reg); err != nil {
		return false
	}else{
		return true
	}

}
func NewConsulRegister(name string, consuladress string) *ConsulRegister {
	return &ConsulRegister{
		Address:                        consuladress, //consul address
		Name:                           name,
		Tag:                            []string{},
		Port:                           80,
		DeregisterCriticalServiceAfter: time.Duration(1) * time.Minute,
		Interval:                       time.Duration(10) * time.Second,
	}
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
// CONSUL //

// ARCHIVOS //
func getFolder(num int) string {

	var c1 int = num / 1000000
	var n1 int = num - c1 * 1000000

	var c2 int = n1 / 10000
	n1 = n1 - c2 * 10000

	var c3 int = n1 / 100
	var c4 int = n1 % 100

	//fmt.Printf("num[%v] c1[%v] c2[%v]", num, c1, c2)
	return strconv.Itoa(c1)+"/"+strconv.Itoa(c2)+"/"+strconv.Itoa(c3)+"/"+strconv.Itoa(c4)
}
func escribirArchivos(){

	d1 := []byte("{\"Id\":1,\"Data\":{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}}")

	x := make([]int, 10000)
	c := 0
	time1 := time.Now()

	for j, _ := range x {

		//j = j + 1863100
		v := 100
		folder := getFolder(j)
		//cant := uint64(v)

		newpath := filepath.Join("/var/tmp/utils/filtros", folder)
		err := os.MkdirAll(newpath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			fmt.Println("FOLDER ERROR: ", err)
		}

		//time1 := time.Now()
		for i := 0; i < v; i++ {
			err := os.WriteFile("/var/tmp/utils/filtros/"+folder+"/"+strconv.Itoa(i), d1, 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
		//elapsed1 := uint64(time.Since(time1) / time.Nanosecond) / cant
		//fmt.Printf("utils/filtros/%v [%v] [%v]\n", folder, j, elapsed1)
		c++
	}
	elapsed1 := time.Since(time1)
	fmt.Printf("Cantidad %v / Tiempo: [%v]\n", c, elapsed1)

}
func leerArchivos(){
	
	time1 := time.Now()
	i := 0
	for i < 2000 {

		n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
		folder := getFolder(int(n.Int64()))
		file, err := os.Open("/home/admin/Go/pruebas/utils/filtros/"+folder)
		if err != nil{
			fmt.Println(err)
		}
		file.Close()
		byteValue, _ := ioutil.ReadAll(file)
		read(byteValue)
		i++

	}

	elapsed1 := uint64(time.Since(time1) / time.Nanosecond) / 2000
	fmt.Printf("DuracionLectura [%v]", elapsed1)

}
func read(x []byte){
	//
}
// ARCHIVOS //

// TIME LAPSED //
func printelaped(start time.Time, str string){
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
}
// TIME LAPSED //

// MONITORING //
func NewStats() *stats {
	s := stats{}
	s.startTime = time.Now()
	return &s
}
func (s *stats) PrintStats() {
	up, err := time.ParseDuration(fmt.Sprintf("%fs", s.SysUptime.Uptime))
	upstring := "SysUptime Error"
	if err == nil {
		updays := up.Hours() / 24
		switch {
		case updays >= 365:
			upstring = fmt.Sprintf("%.0f years", updays/365)
		case updays >= 1:
			upstring = fmt.Sprintf("%.0f days", updays)
		default: // less than a day
			upstring = up.String()
		}
	}

	fmt.Println("*********************************************************")
	fmt.Printf("go-top - %s  up %s,\t\tload average: %.2f, %.2f, %.2f\n",
		s.LoadAverage.Time.Format("15:04:05"), upstring, s.LoadAverage.One, s.LoadAverage.Five, s.LoadAverage.Fifteen)

	fmt.Printf("Cpu(s): %.1f%%us, %.1f%%sy, %.1f%%ni, %.1f%%id, %.1f%%wa, %.1f%%hi, %.1f%%si, %.1f%%st %.1f%%gu\n",
		s.SysCPUAvg.UserPct, s.SysCPUAvg.SystemPct, s.SysCPUAvg.NicePct, s.SysCPUAvg.IdlePct,
		s.SysCPUAvg.IowaitPct, s.SysCPUAvg.IrqPct, s.SysCPUAvg.SoftIrqPct, s.SysCPUAvg.StealPct,
		s.SysCPUAvg.GuestPct)

	fmt.Printf("Mem:  %9dk total, %9dk used, %9dk free, %9dk buffers\n", s.SysMemK.MemTotal,
		s.SysMemK.MemUsed, s.SysMemK.MemFree, s.SysMemK.Buffers)
	fmt.Printf("Swap: %9dk total, %9dk used, %9dk free, %9dk cached\n", s.SysMemK.SwapTotal,
		s.SysMemK.SwapUsed, s.SysMemK.SwapFree, s.SysMemK.Cached)

	fmt.Println("************************************************************")
	if s.ProcCPUAvg.PossiblePct > 0 {
		cpuHelpText := "[see -help flag to change %cpu]"
		if *coresToPegPtr > 0 {
			cpuHelpText = ""
		}
		fmt.Printf("ProcessName\tRES(k)\t%%CPU\t%%CCPU\t%%MEM\n")
		fmt.Printf("this-process\t%d\t%3.1f\t%2.1f\t%3.1f\t%s\n",
			s.CurProcCPUSample.ProcMemUsedK,
			s.ProcCPUAvg.TotalPct,
			100*s.CurProcCPUSample.Total/s.ProcUptime/float64(1),
			100*float64(s.CurProcCPUSample.ProcMemUsedK)/float64(s.SysMemK.MemTotal),
			cpuHelpText)
		fmt.Println("%CCPU is cumulative CPU usage over this process' life.")
		fmt.Printf("Max this-process CPU possible: %3.f%%\n", s.ProcCPUAvg.PossiblePct)
	}
}
func (s *stats) GatherStats() {
	s.SysUptime = systemstat.GetUptime()
	s.ProcUptime = time.Since(s.startTime).Seconds()

	s.SysMemK = systemstat.GetMemSample()
	s.LoadAverage = systemstat.GetLoadAvgSample()

	s.LastCPUSample = s.CurCPUSample
	s.CurCPUSample = systemstat.GetCPUSample()

	if s.sysCPUSampled { // we need 2 samples to get an average
		s.SysCPUAvg = systemstat.GetCPUAverage(s.LastCPUSample, s.CurCPUSample)
	}
	// we have at least one sample, subsequent rounds will give us an average
	s.sysCPUSampled = true

	s.ProcMemUsedPct = 100 * float64(s.CurProcCPUSample.ProcMemUsedK) / float64(s.SysMemK.MemTotal)

	s.LastProcCPUSample = s.CurProcCPUSample
	s.CurProcCPUSample = systemstat.GetProcCPUSample()
	if s.procCPUSampled {
		s.ProcCPUAvg = systemstat.GetProcCPUAverage(s.LastProcCPUSample, s.CurProcCPUSample, s.ProcUptime)
	}
	s.procCPUSampled = true
}
// MONITORING //

func DirSize(path string) (float64, error) {
    var size int64
    err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            size += info.Size()
        }
        return err
    })
	sizeMB := float64(size) / 1024.0 / 1024.0
    return sizeMB, err
}
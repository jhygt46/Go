package main

import (
	"os"
	"fmt"
	"net"
	"log"
	"time"
	"strings"
	"strconv"
	"runtime"
	"math/big"
	"net/http"
	"io/ioutil"
	"crypto/rand"
	"path/filepath"
	"encoding/json"
    "github.com/valyala/fasthttp"
	"github.com/hashicorp/consul/api"
)

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
	AutoCache bool `json:"AutoCache"` // 0 AUTOMATICO - 1 LISTA CACHE
	ListaCache []int64 `json:"ListaCache"`
	TotalCache int32 `json:"TotalCache"`
}
type ConsulRegister struct {
	Address                        string
	Name                           string
	Tag                            []string
	Port                           int
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
}
type Config struct {
	Fecha time.Time `json:"Fecha"`

	// INIT CACHE
	AutoCache bool `json:"Cachetipo"`
	TotalCache int32 `json:"TotalCache"`
	CountCache int32 `json:"CountCache"`

	// TOTAL CACHE
	MetricTime time.Time `json:"MetricTime"`
	MetricCount int64 `json:"MetricCount"`

	// START CACHE/FILES
	MetricStart bool `json:"MetricStart"`
	MetricCache int64 `json:"MetricCache"`
	MetricFile int64 `json:"MetricFile"`
}
type MyHandler struct {
	Minicache map[uint64]*Data
	Config Config
}
type PostResponse struct {
	Consulname string `json:"Consulname"`
	Consulhost string `json:"Consulip"`
	AutoCache bool `json:"AutoCache"` // 0 AUTOMATICO - 1 LISTA CACHE
	ListaCache []int64 `json:"ListaCache"`
	TotalCache int32 `json:"TotalCache"`
}
type PostRequest struct {
	Id string `json:"Id"`
	Ip string `json:"Ip"`
	Init bool `json:"Init"`
	Consul bool `json:"Consul"`
	Time time.Time `json:"Time"`
}
type StatusServer struct {
	Cpu statusCpu `json:"Cpu"`
	SizeMb float64 `json:"DiskMb"`
	Memory statusMemory `json:"Memory"`
}
type statusCpu struct {
	CpuUsage float64 `json:"CpuUsage"`
	IdleTicks float64 `json:"IdleTicks"`
	TotalTicks float64 `json:"TotalTicks"`
}
type statusMemory struct {
	Alloc uint64 `json:"Alloc"`
	TotalAlloc uint64 `json:"TotalAlloc"`
	Sys uint64 `json:"Sys"`
	NumGC uint32 `json:"NumGC"`
}
func main() {

	h := &MyHandler{}
	h.initServer()
	h.statusServer()
	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)
	
}
func (h *MyHandler) initServer() {
	
	h.Config.CountCache = 0
	
	params := PostRequest{}
	//params.Id = getInstanceId()
	params.Id = "ami-636388377355";
	params.Ip = LocalIP()
	params.Init = true
	params.Time = time.Now()
	res := postRequest("http://localhost/init", params)
	
	h.Minicache = make(map[uint64]*Data, res.TotalCache)
	h.Config.TotalCache = res.TotalCache

	if !res.AutoCache {
		for _, v := range res.ListaCache {
			jsonFiltro, err := os.Open("/var/filtros/"+getFolder64(uint64(v)))
			if err == nil {
				byteValueFiltro, _ := ioutil.ReadAll(jsonFiltro)
				data := Data{}
				if err := json.Unmarshal(byteValueFiltro, &data); err == nil {
					h.Minicache[uint64(v)] = &data
					h.Config.CountCache++
				}
			}
			defer jsonFiltro.Close()
		}
		if h.Config.CountCache >= h.Config.TotalCache {
			h.Config.AutoCache = false
		}
	}else{
		h.Config.AutoCache = true
	}

	params.Init = false
	if consulRegister(res.Consulname, res.Consulhost) {
		params.Consul = true
		params.Time = time.Now()
		postRequest("http://localhost/init", params)
	}else{
		params.Consul = false
		params.Time = time.Now()
		postRequest("http://localhost/init", params)
	}
	h.Config.Fecha = time.Now()

}
func (h *MyHandler) statusServer() {
	
	h.Config.CountCache = 0
	
	params := PostRequest{}
	//params.Id = getInstanceId()
	params.Id = "ami-636388377355";
	params.Ip = LocalIP()
	params.Init = true
	params.Time = time.Now()
	res := postRequest("http://localhost/init", params)
	
	h.Minicache = make(map[uint64]*Data, res.TotalCache)
	h.Config.TotalCache = res.TotalCache

	if !res.AutoCache {
		for _, v := range res.ListaCache {
			jsonFiltro, err := os.Open("/var/filtros/"+getFolder64(uint64(v)))
			if err == nil {
				byteValueFiltro, _ := ioutil.ReadAll(jsonFiltro)
				data := Data{}
				if err := json.Unmarshal(byteValueFiltro, &data); err == nil {
					h.Minicache[uint64(v)] = &data
					h.Config.CountCache++
				}
			}
			defer jsonFiltro.Close()
		}
		if h.Config.CountCache >= h.Config.TotalCache {
			h.Config.AutoCache = false
		}
	}else{
		h.Config.AutoCache = true
	}

	params.Init = false
	if consulRegister(res.Consulname, res.Consulhost) {
		params.Consul = true
		params.Time = time.Now()
		postRequest("http://localhost/init", params)
	}else{
		params.Consul = false
		params.Time = time.Now()
		postRequest("http://localhost/init", params)
	}
	h.Config.Fecha = time.Now()

}
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	switch string(ctx.Path()) {
	case "/filtro":

		ctx.Response.Header.Set("Content-Type", "application/json")
		id := read_int64(ctx.QueryArgs().Peek("id")) 
		if res, found := h.Minicache[id]; found {
			json.NewEncoder(ctx).Encode(res)
			if h.Config.MetricStart { h.Config.MetricCache++ }
		}else{
			jsonFiltro, err := os.Open("/var/filtros/"+getFolder64(id))
			if err == nil {
				byteValueFiltro, _ := ioutil.ReadAll(jsonFiltro)
				if h.Config.AutoCache {
					data := Data{}
					if err := json.Unmarshal(byteValueFiltro, &data); err == nil {
						h.Minicache[uint64(id)] = &data
						h.Config.CountCache++
						if h.Config.CountCache >= h.Config.TotalCache {
							h.Config.AutoCache = false
						}
					}
				}
				if h.Config.MetricStart { h.Config.MetricFile++ }
				fmt.Fprintf(ctx, string(byteValueFiltro))
			}else{
				ctx.Error("Not Found", fasthttp.StatusNotFound)
			}
			defer jsonFiltro.Close()
		}
		h.Config.MetricCount++

	case "/monitoring":

		StatusServer := StatusServer{}

		cpu := true
		disk := true
		mem := true

		if cpu {
			StatusServer.Cpu = GetMonitoringsCpu()
		}
		if disk {
			SizeMb, err := DirSize("/var/utils")
			if err == nil {
				StatusServer.SizeMb = SizeMb
			}
		}
		if mem {
			StatusServer.Memory = PrintMemUsage()
		}	

		Alloc, TotalAlloc, Sys, NumGC := PrintMemUsage()
		fmt.Printf("Alloc = %v Mib\n TotalAlloc = %v Mib\n Sys = %v Mib\n NumGC = %v\n", Alloc, TotalAlloc, Sys, NumGC)

	case "/health":
		fmt.Fprintf(ctx, "OK");
	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
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
func read_int64(data []byte) uint64 {
    var x uint64
    for _, c := range data {
        x = x * 10 + uint64(c - '0')
    }
    return x
}
func getUrladminResponse(url string) *adminResponse {

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
func postRequest(url string, post PostRequest) *PostResponse {
	
	u, err := json.Marshal(post)
	if err != nil {
		panic(err)
	}
	req := fasthttp.AcquireRequest()
	req.SetBody(u)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(url)
	res := fasthttp.AcquireResponse()
	var resp PostResponse
	if err := fasthttp.Do(req, res); err == nil {
		defer fasthttp.ReleaseRequest(req)
		body := res.Body()
		json.Unmarshal(body, &resp)
		defer fasthttp.ReleaseResponse(res)
	}
	return &resp

}
func statusRequest(url string, post statusParams) *statusResponse {
	
	u, err := json.Marshal(post)
	if err != nil {
		panic(err)
	}
	req := fasthttp.AcquireRequest()
	req.SetBody(u)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(url)
	res := fasthttp.AcquireResponse()
	var resp statusResponse
	if err := fasthttp.Do(req, res); err == nil {
		defer fasthttp.ReleaseRequest(req)
		body := res.Body()
		json.Unmarshal(body, &resp)
		defer fasthttp.ReleaseResponse(res)
	}
	return &resp

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
func getFolder64(num uint64) string {

	c1, n1 := divmod(num, 1000000)
	c2, n2 := divmod(n1, 10000)
	c3, c4 := divmod(n2, 10000)
	return strconv.FormatUint(c1, 10)+"/"+strconv.FormatUint(c2, 10)+"/"+strconv.FormatUint(c3, 10)+"/"+strconv.FormatUint(c4, 10)

}
func divmod(numerator, denominator uint64) (quotient, remainder uint64) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}
func escribirArchivos(path string){

	d1 := []byte("{\"Id\":1,\"Data\":{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}}")

	x := make([]int64, 100000)
	c := 0
	time1 := time.Now()

	for j, _ := range x {

		v := 100
		folder := getFolder64(uint64(j))

		newpath := filepath.Join(path, folder)
		err := os.MkdirAll(newpath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			fmt.Println("FOLDER ERROR: ", err)
		}

		for i := 0; i < v; i++ {
			err := os.WriteFile(path+"/"+folder+"/"+strconv.Itoa(i), d1, 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
		c++

	}
	elapsed1 := time.Since(time1)
	fmt.Printf("Cantidad %v / Tiempo: [%v]\n", c, elapsed1)

}
func leerArchivos(path string){
	
	time1 := time.Now()
	i := 0
	for i < 2000 {

		n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
		folder := getFolder64(n.Uint64())
		file, err := os.Open(path+"/"+folder)
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

// MONITORING //
func getCPUSample() (idle, total uint64) {
    contents, err := ioutil.ReadFile("/proc/stat")
    if err != nil {
        return
    }
    lines := strings.Split(string(contents), "\n")
    for _, line := range(lines) {
        fields := strings.Fields(line)
        if fields[0] == "cpu" {
            numFields := len(fields)
            for i := 1; i < numFields; i++ {
                val, err := strconv.ParseUint(fields[i], 10, 64)
                if err != nil {
                    fmt.Println("Error: ", i, fields[i], err)
                }
                total += val // tally up all the numbers to get total ticks
                if i == 4 {  // idle is the 5th field in the cpu line
                    idle = val
                }
            }
            return
        }
    }
    return
}
func GetMonitoringsCpu() statusCpu {

	statusCpu := statusCpu{}

	idle0, total0 := getCPUSample()
	time.Sleep(3 * time.Second)
	idle1, total1 := getCPUSample()

	statusCpu.IdleTicks = float64(idle1 - idle0)
	statusCpu.TotalTicks = float64(total1 - total0)
	statusCpu.CpuUsage = 100 * (statusCpu.TotalTicks - statusCpu.IdleTicks) / statusCpu.TotalTicks

	return statusCpu

}
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
func PrintMemUsage() statusMemory {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	statusMemory := statusMemory{}
	statusMemory.Alloc = bToMb(m.Alloc)
	statusMemory.TotalAlloc = bToMb(m.TotalAlloc)
	statusMemory.Sys = bToMb(m.Sys)
	statusMemory.NumGC = m.NumGC
	return statusMemory
}
func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
// MONITORING //


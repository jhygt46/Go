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
)

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
	consulname string `json:"consulname"`
	consulip string `json:"consulip"`
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
	h.initServer()

	/*
	time1 := time.Now()
	if consulname, consulhost, err := initServer(); err {
		printelaped(time1, "SERVER INICIADO...")
		time2 := time.Now()
		if consulRegister(consulname, consulhost) {
			printelaped(time2, "CONSUL INICIADO...")
			
			fasthttp.ListenAndServe(":80", h.HandleFastHTTP)
		}else{
			fmt.Println("Consul Register ERROR")
		}
	}else{
		fmt.Println("ERROR AL INICIAR SERVIDOR")
	}
	*/
	
}

func (h *MyHandler) initServer() {

	adminResponse := &adminResponse{}
	err := getUrl("http://18.118.129.19/", adminResponse)
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Println("adminResponse")
	fmt.Println(*adminResponse)
	/*
	id := getInstanceId()
	ip := LocalIP()
	resp, err := http.Get("http://18.118.129.19/init/?id="+id+"&ip="+ip)
	if err != nil {
		log.Fatalln(err)
		return "", "", false
	}
	body, err := ioutil.ReadAll(resp.Body)
   	if err != nil {
		log.Fatalln(err)
		return "", "", false
   	}

	fmt.Println(string(body))

	if string(body) == "OK" {
		return "filtro1", "10.128.0.4:8500", true
	}else{
		return "", "", false
	}
	*/

}

func getUrl(url string, target interface{}) error {

	myClient := &http.Client{Timeout: 10 * time.Second}
    r, err := myClient.Get(url)
    if err != nil {
        return err
    }
    defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	fmt.Println("BODY")
	fmt.Println(string(body))

    return json.NewDecoder(r.Body).Decode(target)

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
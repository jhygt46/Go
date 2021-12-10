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
	"encoding/json"
    "github.com/valyala/fasthttp"
	//"utils/utils"
	//"github.com/aws/aws-sdk-go-v2/aws"
	//"github.com/aws/aws-sdk-go-v2/config"
	//"github.com/aws/aws-sdk-go-v2/service/ec2"
	//"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// SERVICE KUBERNET//

type Kubernet struct {
	Servicios []Servicios
}


type Servicio struct {
	Nombre string `json:"Nombre"`
	Acl_tipo int8 `json:"Acl_tipo"`
	Acl_valor string `json:"Acl_valor"`
	Backends []Backends `json:"Backends"`
}
type Backends struct {
	Active bool `json:"Active"`
	Fecha time.Time `json:"Fecha"`
	Lista_backends []Backend `json:"Lista_backends"`
}
type Backend struct {
	Acls []Acl `json:"Acls"`
	Servers []Server `json:"Servers"`
	Backend string`json:"Backend"`
	Consulname string`json:"Consulname"`
}
type Acl struct {
	Nombre string`json:"Nombre"`
	Param string`json:"Param"`
	Tipo int8 `json:"Tipo"`
	Valor1 int64 `json:"Valor1"`
	Valor2 int64 `json:"Valor2"`
}
type Server struct {
	Ip string `json:"Ip"`
	Id string `json:"Id"`
	Cpu Cpu `json:"Cpu"`
	Memory Memory `json:"Memory"`
	DiskMb int32 `json:"Disk"`
}
type Cpu struct {
	Param1 float64 `json:"Param1"`
}
type Memory struct {
	Param1 float64 `json:"Param1"`
}
// SERVICE KUBERNET//



/*
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
type Daemon struct {
	Servicios []Servicio `json:"Servicios"`
}
type adminResponse struct {
	Consulname string `json:"Consulname"`
	Consulhost string `json:"Consulip"`
	Cachetipo int8 `json:"Cachetipo"` // 0 AUTOMATICO - 1 LISTA CACHE
	ListaCache []int64 `json:"ListaCache"`
	TotalCache int32 `json:"TotalCache"`
}
type EC2API interface {
	CreateImage(ctx context.Context, params *ec2.CreateImageInput, optFns ...func(*ec2.Options)) (*ec2.CreateImageOutput, error)
	RunInstances(ctx context.Context, params *ec2.RunInstancesInput, optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)
	CreateTags(ctx context.Context, params *ec2.CreateTagsInput, optFns ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error)
	DeregisterImage(ctx context.Context, params *ec2.DeregisterImageInput, optFns ...func(*ec2.Options)) (*ec2.DeregisterImageOutput, error)
	TerminateInstances(ctx context.Context, params *ec2.TerminateInstancesInput, optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}
*/
// TYPES //

type PostRequest struct {
	Id string `json:"Id"`
	Ip string `json:"Ip"`
	Init bool `json:"Init"`
	Consul bool `json:"Consul"`
	Time time.Time `json:"Time"`
}

type Config struct {
	Tiempo time.Duration `json:"Tiempo"`
}
type MyHandler struct {
	 
	Servicios []*Servicio `json:"Servicios"`
}

func main() {

	acl1 := Acl{ Nombre: "Acl1", Param: "Param1", Tipo: 1, Valor1: 0, Valor2: 100 }
	ser1 := Server{ Ip: "127.0.0.1", Id: "ami-664875373826", Cpu: Cpu{ Param1: 100 }, Memory: Memory{ Param1: 105 }, DiskMb: 1000 }

	bck1 := Backend{ Backend: "BACKNAME", Consulname: "BACKCONSUL", Acls: []Acl{acl1}, Servers: []Server{ser1} }

	l_bck1 := Backends{ Active: true, Fecha: time.Now(), Lista_backends: []Backend{bck1} }
	
	Servicios := make([]*Servicio, 0)
	Servicios = append(Servicios, &Servicio{ Nombre: "Filtro", Acl_tipo: 1, Acl_valor: "/filtro/", Backends: []Backends{l_bck1}  })
	
	//dae := readFile("daemon.json")
	pass := &MyHandler{ Conf: &Config{}, Servicios: Servicios }

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
		fasthttp.ListenAndServe(":81", pass.HandleFastHTTP)
	}()
	
	if err := run(con, pass, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
    ctx.Response.Header.Set("Access-Control-Allow-Headers", "authorization, content-type, set-cookie, cookie, server")
    ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
    ctx.Response.Header.Set("Access-Control-Allow-Credentials", "false")

	if string(ctx.Method()) == "POST" {
		params := ctx.PostBody()
		fmt.Println(string(params))
		switch string(ctx.Path()) {
		case "/init":
			//fmt.Println(ctx.RemoteAddr())
			var res PostRequest
			if err := json.Unmarshal(params, &res); err == nil {
				fmt.Println("Id", res.Id)
				fmt.Println("Ip", res.Ip)
				fmt.Println("Init", res.Init)
				fmt.Println("Consul", res.Consul)
				fmt.Println("Time", res.Time)
				fmt.Println(res)
			}else{
				fmt.Println(err)
			}
			json.NewEncoder(ctx).Encode(res)
		case "/status":
			var res statusParams
			if err := json.Unmarshal(params, &res); err == nil {
				//fmt.Println("DiskMb: ", res.DiskMb)
				h.statusServer(res)
			}else{
				fmt.Println(err)
			}
		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}
	//ctx.Error("Not Found", fasthttp.StatusNotFound)

}


func (h *MyHandler) statusServer(params statusParams){
 
	for _, bck := range h.Servicios[param.serv].Backends {
		if bck.Active {
			for _, l_bck := range bck.Lista_backends {
				if l_bck
			}
		}
	}

}

// DAEMON //
func (h *MyHandler) StartDaemon() {

	for _, serv := range h.Servicios {
		fmt.Println("SERVICIO NOMBRE", serv.Nombre)
		fmt.Println("SERVICIO TIPO", serv.Acl_tipo)
		fmt.Println("SERVICIO ACL VALOR", serv.Acl_valor)
		for _, bck := range serv.Backends {
			fmt.Println("BACKEND ACTIVE", bck.Active)
			fmt.Println("BACKEND FECHA", bck.Fecha)
			for _, l_bck := range bck.Lista_backends {
				fmt.Println("L_BCK: ", l_bck.Backend)
				fmt.Println("L_BCK: ", l_bck.Consulname)
				for _, acls := range l_bck.Acls {
					fmt.Println("ACL NOMBRE: ", acls.Nombre)
					fmt.Println("ACL PARAM: ", acls.Param)
					fmt.Println("ACL TIPO: ", acls.Tipo)
					fmt.Println("ACL VALOR 1: ", acls.Valor1)
					fmt.Println("ACL VALOR 2: ", acls.Valor2)
				}
				for _, servers := range l_bck.Servers {
					fmt.Println("SERVER ID: ", servers.Id)
					fmt.Println("SERVER IP: ", servers.Ip)
					fmt.Println("SERVER CPU: ", servers.Cpu)
					fmt.Println("SERVER MEMORY: ", servers.Memory)
					fmt.Println("SERVER DISKMB: ", servers.DiskMb)
				}
			}
		}
	}
	h.Conf.Tiempo = 20 * time.Second

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

/*
// AWS USED FUNCTION //
func delete_image(ami string) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.DeregisterImageInput{
		ImageId: aws.String(ami),
	}

	resp, err := DelImage(context.TODO(), client, input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)

}
func terminate_instance() {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)


	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{"i-080a8dad14046e40c"},
	}

	resp, err := DelInstance(context.TODO(), client, input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(*resp)

}
func create_instance(ImageId string, TagName string, TagValue string) string {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	min := int32(1)
	max := int32(1)

	SecurityGroupId := []string{"sg-0dbcca3589e78cefd"}

	input := &ec2.RunInstancesInput{
		ImageId:      &ImageId,
		InstanceType: types.InstanceTypeT2Nano,
		MinCount:     &min,
		MaxCount:     &max,
		SecurityGroupIds: SecurityGroupId,
	}

	result, err := MakeInstance(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error creating an instance:")
		fmt.Println(err)
		return ""
	}

	tagInput := &ec2.CreateTagsInput{
		Resources: []string{*result.Instances[0].InstanceId},
		Tags: []types.Tag{
			{
				Key:   &TagName,
				Value: &TagValue,
			},
		},
	}

	_, err = MakeTags(context.TODO(), client, tagInput)
	if err != nil {
		fmt.Println("Got an error tagging the instance:")
		fmt.Println(err)
		return ""
	}

	fmt.Println(result)
	return *result.Instances[0].InstanceId
	
}
func create_image(InstanceId string, Nombre string, Descripcion string) string {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.CreateImageInput{
		Description: &Descripcion,
		InstanceId:  &InstanceId,
		Name:        &Nombre,
		BlockDeviceMappings: []types.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sda1"),
				NoDevice:   aws.String(""),
			},
			{
				DeviceName: aws.String("/dev/sdb"),
				NoDevice:   aws.String(""),
			},
			{
				DeviceName: aws.String("/dev/sdc"),
				NoDevice:   aws.String(""),
			},
		},
	}

	resp, err := MakeImage(context.TODO(), client, input)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return *resp.ImageId

}
// AWS USED FUNCTION //

// AWS NATIVE FUNCTIONS //
func DelImage(c context.Context, api EC2API, input *ec2.DeregisterImageInput) (*ec2.DeregisterImageOutput, error) {
	return api.DeregisterImage(c, input)
}
func DelInstance(c context.Context, api EC2API, input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	return api.TerminateInstances(c, input)
}
func MakeInstance(c context.Context, api EC2API, input *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error) {
	return api.RunInstances(c, input)
}
func MakeTags(c context.Context, api EC2API, input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return api.CreateTags(c, input)
}
func MakeImage(c context.Context, api EC2API, input *ec2.CreateImageInput) (*ec2.CreateImageOutput, error) {
	return api.CreateImage(c, input)
}
// AWS FUNCTIONS //
*/

// UTILS //

/*

*/

// UTILS //
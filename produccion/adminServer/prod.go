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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// TYPES //
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
type EC2API interface {
	CreateImage(ctx context.Context, params *ec2.CreateImageInput, optFns ...func(*ec2.Options)) (*ec2.CreateImageOutput, error)
	RunInstances(ctx context.Context, params *ec2.RunInstancesInput, optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)
	CreateTags(ctx context.Context, params *ec2.CreateTagsInput, optFns ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error)
	DeregisterImage(ctx context.Context, params *ec2.DeregisterImageInput, optFns ...func(*ec2.Options)) (*ec2.DeregisterImageOutput, error)
	TerminateInstances(ctx context.Context, params *ec2.TerminateInstancesInput, optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}
type adminResponse struct {
	consulname string `json:"consulname"`
	consulip string `json:"consulip"`
}
// TYPES //

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

	consul := adminResponse{ consulname: "filtro1", consulip: "10.128.0.4:8500" }
	fmt.Println(consul)
	//fmt.Println(h.Conf)
	//fmt.Println(*h.Dae)
	ctx.Response.Header.Set("Content-Type", "application/json")
	json.NewEncoder(ctx).Encode(consul)

}


// DAEMON //
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
// DAEMON //


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


// UTILS //
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
// UTILS //
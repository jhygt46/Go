package main

import (
	"os"
	"io"
	"log"
	"fmt"
	"flag"
	"math"
	"time"
	//"bytes"
	"bufio"
	"errors"
	"syscall"
	"context"
	//"io/ioutil"
	"os/signal"
	//"archive/tar"
	"encoding/json"
	"path/filepath"
    "github.com/valyala/fasthttp"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mitchellh/go-homedir"
    "github.com/docker/docker/pkg/archive"
)
type Config struct {
	Id int8 `json:"Id"`
	Fecha time.Time `json:"Fecha"`
	Tiempo time.Duration `json:"Time"`
}
type MyHandler struct {
	Conf *Config
	cli *client.Client
}
type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}
type ErrorDetail struct {
	Message string `json:"message"`
}

var dockerRegistryUserID = "111"

func main() {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	
	if imageBuild("/var/docker-images/filtros/Dockerfile", cli) {
		fmt.Println("IMAGEN CREADA")
	}else{
		fmt.Println("ERROR CREAR IMAGEN")
	}
	
	pass := &MyHandler{ Conf: &Config{ Id: 8, Fecha: time.Now() }, cli: cli }

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
func GetContext(filePath string) io.Reader {

    filePaths, _ := homedir.Expand(filePath)
    ctx, _ := archive.TarWithOptions(filePaths, &archive.TarOptions{})
    return ctx
}
func imageBuild(s string, cli *client.Client) bool {

	os.Chdir("/var/dockers-images/filtros")

	ctx := context.Background()
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	for _, image := range images {
		for _, img := range image.RepoTags{
			if img == "filtrogo:latest" {
				fmt.Println("image encontrada")
				dels, err := cli.ImageRemove(context.Background(), image.ID, types.ImageRemoveOptions{Force: true, PruneChildren: false})
				if err != nil {
					panic(err)
				}else{
					fmt.Println("image eliminado")
					fmt.Println(dels)
				}
				
			}
		}
		//fmt.Println(image.ID)
		//fmt.Println(image.Size)
		//fmt.Println(image.VirtualSize)
	}

	buildOptions := types.ImageBuildOptions{
		Tags:   []string{"filtrogo"},
	}

	tar, err := archive.TarWithOptions("/var/docker-images/filtros/", &archive.TarOptions{})
	if err != nil {
		panic(err)
	}


	imageBuildResponse, err := cli.ImageBuild(context.Background(), tar, buildOptions)
	if err != nil {
        log.Fatalf("build error - %s", err)
    }
	io.Copy(os.Stdout, imageBuildResponse.Body)
    defer imageBuildResponse.Body.Close()

	return true

}
func (h *MyHandler) StartDaemon() {

	h.Conf.Tiempo = 5 * time.Second

	errs := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		return nil
	})
	if errs != nil { panic(errs) }

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

	switch string(ctx.Path()) {
	case "/filtro":
		fmt.Fprintf(ctx, "OK");
	case "/health":
		fmt.Fprintf(ctx, "OK");
	case "/info":
		fmt.Fprintf(ctx, "OK");
	case "/metrica":
	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}
	
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
func print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

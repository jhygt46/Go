package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"flag"
	"time"
	"strings"
	"syscall"
	"context"
	"os/signal"
)

// TYPES //

type adminResponse struct {
	Consulname string `json:"Consulname"`
	Consulhost string `json:"Consulip"`
	Cachetipo int8 `json:"Cachetipo"` // 0 AUTOMATICO - 1 LISTA CACHE
	ListaCache []int64 `json:"ListaCache"`
	TotalCache int32 `json:"TotalCache"`
}
type Config struct {
	Tiempo time.Duration `json:"Tiempo"`
}
type PostRequest struct {
	Id string `json:"Id"`
	Ip string `json:"Ip"`
	Init bool `json:"Init"`
	Consul bool `json:"Consul"`
	Time time.Time `json:"Time"`
}
type MyHandler struct {
	Conf *Config `json:"Conf"`
	Servicios []*Servicio `json:"Servicios"`
}
type Servicio struct {
	Backend []Backend `json:"Backend"`
	Nombre string`json:"Nombre"`
	Valor string`json:"Valor"`
	Tipo int8 `json:"Tipo"`
	Activo bool `json:"Activo"`
}
type Backend struct {
	Acls []Acl `json:"Acls"`
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

func main() {

	Servicios := make([]*Servicio, 0)

	acl1 := Acl{ Nombre: "alpha", Tipo: 1, Valor1: 100, Valor2: 200, Param: "id" }
	acl2 := Acl{ Nombre: "alpha", Tipo: 2, Valor1: 100, Valor2: 200, Param: "id" }

	acl3 := Acl{ Nombre: "alpha", Tipo: 1, Valor1: 201, Valor2: 300, Param: "id" }
	acl4 := Acl{ Nombre: "alpha", Tipo: 2, Valor1: 201, Valor2: 300, Param: "id" }

	back1 := Backend{ Consulname: "filtro", Backend: "Filtro", Acls: []Acl{ acl1, acl2 } }
	back2 := Backend{ Consulname: "filtro", Backend: "Filtro", Acls: []Acl{ acl3, acl4 } }

	acl5 := Acl{ Nombre: "alpha", Tipo: 1, Valor1: 100, Valor2: 200, Param: "id" }
	acl6 := Acl{ Nombre: "alpha", Tipo: 2, Valor1: 100, Valor2: 200, Param: "id" }

	acl7 := Acl{ Nombre: "alpha", Tipo: 1, Valor1: 201, Valor2: 300, Param: "id" }
	acl8 := Acl{ Nombre: "alpha", Tipo: 2, Valor1: 201, Valor2: 300, Param: "id" }

	back3 := Backend{ Consulname: "auto", Backend: "Auto", Acls: []Acl{ acl5, acl6 } }
	back4 := Backend{ Consulname: "auto", Backend: "Auto", Acls: []Acl{ acl7, acl8 } }

	Servicios = append(Servicios, &Servicio{ Backend: []Backend{ back1, back2 }, Nombre: "is_filtro", Tipo: 1, Valor: "/filtros/", Activo: true })
	Servicios = append(Servicios, &Servicio{ Backend: []Backend{ back3, back4 }, Nombre: "is_auto", Tipo: 1, Valor: "/auto/", Activo: false })
	
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
	
	if err := run(con, pass, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}

func (h *MyHandler) create_config_file() string{

	var b strings.Builder
	fmt.Fprintf(&b, "global\n\tlog /dev/log    local0\n\tlog /dev/log    local1 notice\n\tchroot /var/lib/haproxy\n\tstats socket ipv4@127.0.0.1:9999 level admin\n\tstats socket /var/run/haproxy.sock mode 666 level admin\n\tstats timeout 2m\n\tuser haproxy\n\tgroup haproxy\n\tdaemon\n\n\tca-base /etc/ssl/certs\n\tcrt-base /etc/ssl/private\n\n\tssl-default-bind-ciphers ECDH+AESGCM:DH+AESGCM:ECDH+AES256:DH+AES256:ECDH+AES128:DH+AES:RSA+AESGCM:RSA+AES:!aNULL:!MD5:!DSS\n\tssl-default-bind-options no-sslv3\n\n")
	fmt.Fprintf(&b, "defaults\n\tlog     global\n\tmode    http\n\toption  httplog\n\toption  dontlognull\n\ttimeout connect 5000\n\ttimeout client  50000\n\ttimeout server  50000\n\terrorfile 400 /etc/haproxy/errors/400.http\n\terrorfile 403 /etc/haproxy/errors/403.http\n\terrorfile 408 /etc/haproxy/errors/408.http\n\terrorfile 500 /etc/haproxy/errors/500.http\n\terrorfile 502 /etc/haproxy/errors/502.http\n\terrorfile 503 /etc/haproxy/errors/503.http\n\terrorfile 504 /etc/haproxy/errors/504.http\n\n")
	fmt.Fprintf(&b, "frontend apache_front\n\n\tbind *:80\n\n")

	for _, servicio := range h.Servicios {
		if servicio.Activo {
			switch servicio.Tipo {
			case 1:
				fmt.Fprintf(&b, "\tacl %s path_beg %s\n", servicio.Nombre, servicio.Valor)
			default:
				fmt.Fprintf(&b, "")
			}
			for i, backend := range servicio.Backend {
				for j, acl := range backend.Acls {
					fmt.Fprintf(&b, "\tacl %s%d%d ", acl.Nombre, i, j)
					switch acl.Tipo {
					case 1:
						fmt.Fprintf(&b, "urlp_val(%s) %d:%d\n", acl.Param, acl.Valor1, acl.Valor2)
					case 2:
						fmt.Fprintf(&b, "urlp_reg(%s) ^[%d-%d]\n", acl.Param, acl.Valor1, acl.Valor2)
					default:
						fmt.Fprintf(&b, "")
					}
				}
			}
			fmt.Fprintf(&b, "\n")
		}
	}

	for _, servicio := range h.Servicios {
		if servicio.Activo {
			for i, backend := range servicio.Backend {
				fmt.Fprintf(&b, "\tuse_backend %s%d if %s", backend.Backend, i, servicio.Nombre)
				for j, acl := range backend.Acls {
					fmt.Fprintf(&b, " %s%d%d", acl.Nombre, i, j)
				}
				fmt.Fprintf(&b, "\n")
			}
		}
	}

	fmt.Fprintf(&b, "\n")

	for _, servicio := range h.Servicios {
		if servicio.Activo {
			for i, backend := range servicio.Backend {
				fmt.Fprintf(&b, "backend %s%d\n\tbalance roundrobin\n\tserver-template mywebapp 10 _%s%d._tcp.service.consul resolvers consul resolve-opts allow-dup-ip resolve-prefer ipv4 check\n\ttimeout connect 1m\n\ttimeout server 1m\n\n", backend.Backend, i, backend.Consulname, i)
			}
		}
	}
	
	ip := "10.128.0.4"
	port := "8600"
	fmt.Fprintf(&b, "resolvers consul\n\tnameserver consul %s:%s\n\taccepted_payload_size 8192\n\thold valid 5s\n\n", ip, port)

	return b.String()

}


// DAEMON //
func (h *MyHandler) StartDaemon() {

	fmt.Println("DAEMON: ", h.Conf.Tiempo)
	h.Conf.Tiempo = 20 * time.Second

	f, err := os.Create("harpoxy.cfg")
	if err == nil {
		n, _ := f.WriteString(h.create_config_file())
		fmt.Printf("bytes: %d", n)
	}
	defer f.Close()

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



// UTILS //
/*
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
*/
// UTILS //
package main

import (
    "net"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
)

func NewConsulRegister() *ConsulRegister {
	return &ConsulRegister{
		Address:                        10.128.0.4:8500, //consul address
		Name:                           "unknown",
		Tag:                            []string{},
		Port:                           80,
		DeregisterCriticalServiceAfter: time.Duration(1) * time.Minute,
		Interval:                       time.Duration(10) * time.Second,
	}
}
 
// ConsulRegister consul service register
type ConsulRegister struct {
	Address                        string
	Name                           string
	Tag                            []string
	Port                           int
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
}

func main() {

	r := NewConsulRegister()
	config := api.DefaultConfig()
	config.Address = r.Address
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	agent := client.Agent()

	IP := LocalIP()
	reg := &api.AgentServiceRegistration{
		 ID: fmt.Sprintf("%v-%v-%v", r.Name, IP, r.Port), // Name of the service node
		 Name: r.Name, // service name
		 Tags: r.Tag, // tag, can be empty
		 Port: r.Port, // service port
		 Address: IP, // Service IP
		 Check: &api.AgentServiceCheck{ // Health Check
			 Interval: r.Interval.String(), // Health check interval
			 GRPC: fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Name), // grpc support, address to perform health check, service will be passed to Health.Check function
			 DeregisterCriticalServiceAfter: r.DeregisterCriticalServiceAfter.String(), // Deregistration time, equivalent to expiration time
		},
	}
 
	if err := agent.ServiceRegister(reg); err != nil {
		return err
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
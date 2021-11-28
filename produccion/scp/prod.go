package main

import (
	"fmt"
	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
	"os"
)

func main() {
	// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library
	clientConfig, _ := auth.PrivateKey("root", "/root/.ssh/id_rsa", ssh.InsecureIgnoreHostKey())

	client := scp.NewClient("3.142.90.232:22", &clientConfig)

	err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return
	}

	f, errs := os.OpenFile("/var/hola.txt", os.O_RDONLY|os.O_CREATE, 0666)
	if errs != nil {
		fmt.Println("Error open file ", errs)
	}

	defer client.Close()
	defer f.Close()

	err = client.CopyFromRemote(f, "/root/hola.txt")
	//err = client.CopyFile(f, "/root/test.txt", "0655")
	if err != nil {
		fmt.Println("Error while copying file ", err)
	}

	

}
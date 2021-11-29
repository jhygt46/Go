package main

import (
	"fmt"
	"io/ioutil"
	"github.com/povsister/scp"
)

func main() {
	
	privPEM, err := ioutil.ReadFile("/root/.ssh/id_rsa")
	sshConf, err := scp.NewSSHConfigFromPrivateKey("root", privPEM, "buenanelson")
	scpClient, err := scp.NewClient("18.117.117.108:22", sshConf, &scp.ClientOption{})
	defer scpClient.Close()

	err = scpClient.CopyFileFromRemote("/var/dd.txt", "/var/dd.txt", &scp.FileTransferOption{})
	if err != nil {
		fmt.Println(err)
	}

	err := scpClient.CopyDirToRemote("/var/copy", "/var/copy", &scp.DirTransferOption{})
	if err != nil {
		fmt.Println(err)
	}

}
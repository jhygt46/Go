package main

import (
	"fmt"
	"time"
	"io/ioutil"
	"github.com/povsister/scp"
)

func main() {
	
	privPEM, err := ioutil.ReadFile("/root/.ssh/id_rsa")
	if err != nil {
		fmt.Println("errx")
		fmt.Println(err)
	}
	sshConf, err := scp.NewSSHConfigFromPrivateKey("root", privPEM, "buenanelson")
	if err != nil {
		fmt.Println("erry")
		fmt.Println(err)
	}
	scpClient, err := scp.NewClient("18.117.117.108:22", sshConf, &scp.ClientOption{})
	if err != nil {
		fmt.Println("errz")
		fmt.Println(err)
	}
	defer scpClient.Close()

	now := time.Now()
	err1 := scpClient.CopyFileFromRemote("/var/dd.txt", "/var/dd.txt", &scp.FileTransferOption{})
	if err1 != nil {
		fmt.Println("err1")
		fmt.Println(err1)
	}else{
		printelaped(now, "COPY DB")
	}
	err2 := scpClient.CopyDirFromRemote("/var/copy", "/var/copy", &scp.DirTransferOption{})
	if err2 != nil {
		fmt.Println("err2")
		fmt.Println(err2)
	}

}
func printelaped(start time.Time, str string) {
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
	//return time.Now()
}
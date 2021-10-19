package main

import (
	"log"
	"fmt"
	"os/exec"
	"encoding/json"
)

var image struct {
	ImageId string `json:"ImageId"`
}

func main(){
	ExampleCmd_StdoutPipe2()
}

func ExampleCmd_StdoutPipe() {
	cmd := exec.Command("bash", "-c", "aws ec2 create-image --instance-id i-07f96abb2dd303e22 --name Mys1 --description Ans1")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCmd_StdoutPipe2() {
	cmd := exec.Command("bash", "-c", "aws ec2 create-image --instance-id i-07f96abb2dd303e22 --name Mys1 --description Ans1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(stdout).Decode(&image); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ImageId: %s\n", image.ImageId)
}

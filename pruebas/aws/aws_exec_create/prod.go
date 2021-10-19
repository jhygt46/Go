package main

import (
	"log"
	//"fmt"
	"os/exec"
	//"encoding/json"
)

var image struct {
	ImageId string `json:"ImageId"`
}

func main(){
	ExampleCmd_StdoutPipe()
}

func ExampleCmd_StdoutPipe() {
	cmd := exec.Command("bash", "-c", "aws ec2 create-image --instance-id i-0f1afaf7e9156a147 --name 'My server' --description 'An AMI for my server'")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}
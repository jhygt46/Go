package main

import (
	"log"
	//"fmt"
	"time"
	"os/exec"
	"context"
	//"encoding/json"
)

var image struct {
	ImageId string `json:"ImageId"`
}

func main(){
	ExampleCmd_StdoutPipe()
	ExampleCommandContext()
}

func ExampleCmd_StdoutPipe() {
	cmd := exec.Command("bash", "-c", "aws ec2 create-image --instance-id i-07f96abb2dd303e22 --name Mys1 --description Ans1")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCommandContext() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if err := exec.CommandContext(ctx, "aws ec2 create-image --instance-id i-07f96abb2dd303e22 --name Mys2 --description Ans2").Run(); err != nil {
		// This will fail after 100 milliseconds. The 5 second sleep
		// will be interrupted.
	}
}
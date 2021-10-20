package main

import (
	"log"
	"fmt"
	"time"
	"os/exec"
	"encoding/json"
)

var image struct {
	ImageId string `json:"ImageId"`
}

type Ec2 struct {
	OwnerId string `json:"OwnerId"`
	ReservationId string `json:"ReservationId"`
	Groups [] Groups `json:"Groups"`
	Instances [] Instances `json:"Instances"`
}
type Groups struct {
	GroupName string `json:"GroupName"`
	GroupId string `json:"GroupId"`
}
type Instances struct {
	Monitoring Monitoring `json:"Monitoring"`
	PublicDnsName string `json:"PublicDnsName"`
	Platform string `json:"Platform"`
	State State `json:"State"`
	EbsOptimized bool `json:"EbsOptimized"`
	LaunchTime time.Time `json:"LaunchTime"`
	InstanceId string `json:"InstanceId"`
	ImageId string `json:"ImageId"`
}
type Monitoring struct {
	PublicDStatensName string `json:"State"`
}
type State struct {
	Code int32 `json:"Code"`
	Name string `json:"Name"`
}



func main(){
	imageId := CreateImageCmd("i-080a8dad14046e40c")
	fmt.Printf("ImageId: %s\n", imageId)
}

func CreateImageCmd(id string) string {

	result := fmt.Sprintf("aws ec2 create-image --instance-id %s --name Mys1 --description Ans1", id)
	//fmt.Println(result)

	cmd := exec.Command( result)
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
	//return image.ImageId
	return result

}

func CreateInstanceCmd(id string) string {
	cmd := exec.Command("bash", "-c", "aws ec2 create-image --instance-id "+id+" --name Mys1 --description Ans1")
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
	return image.ImageId
}

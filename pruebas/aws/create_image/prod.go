package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "fmt"
)

func main() {
    create_image("i-0f1afaf7e9156a147")
}

func create_image(InstanceId string){

	// Load session from shared config
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    // Create EC2 service client
    svc := ec2.New(sess)

    opts := &ec2.CreateImageInput{
        Description: aws.String("ImageTest"),
        InstanceId:  aws.String(InstanceId),
        Name:        aws.String("ImageTest"),
        BlockDeviceMappings: []*ec2.BlockDeviceMapping{
            {
                DeviceName: aws.String("/dev/sda1"),
                NoDevice:    aws.String(""),
            },
            {
                DeviceName: aws.String("/dev/sdb"),
                NoDevice:    aws.String(""),
            },
            {
                DeviceName: aws.String("/dev/sdc"),
                NoDevice:    aws.String(""),
            },
        },
    }
    resp, err := svc.CreateImage(opts)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("ID: ", resp.ImageId)

}
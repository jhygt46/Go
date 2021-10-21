package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "fmt"
)

//https://docs.aws.amazon.com/code-samples/latest/catalog/go-ec2-create_image_no_block_device.go.html

func main() {
    ImageId := create_image("i-080a8dad14046e40c", "Name1", "Desc1")
    fmt.Println(ImageId)
}

func create_image(InstanceId string, Nombre string, Descripcion string) string {

    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    svc := ec2.New(sess)

    opts := &ec2.CreateImageInput{
        Description: aws.String(Descripcion),
        InstanceId:  aws.String(InstanceId),
        Name:        aws.String(Nombre),
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
        return "0"
    }

    return *resp.ImageId

}

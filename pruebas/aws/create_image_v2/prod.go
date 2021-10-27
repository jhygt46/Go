// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2CreateImageAPI interface {
	CreateImage(ctx context.Context, params *ec2.CreateImageInput, optFns ...func(*ec2.Options)) (*ec2.CreateImageOutput, error)
}

func main() {

	imageId := create_image("i-080a8dad14046e40c", "NomImg1", "DescImg1")
	fmt.Println("imageId", imageId)

}

func create_image(InstanceId string, Nombre string, Descripcion string) string {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.CreateImageInput{
		Description: &Descripcion,
		InstanceId:  &InstanceId,
		Name:        &Nombre,
		BlockDeviceMappings: []types.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sda1"),
				NoDevice:   aws.String(""),
			},
			{
				DeviceName: aws.String("/dev/sdb"),
				NoDevice:   aws.String(""),
			},
			{
				DeviceName: aws.String("/dev/sdc"),
				NoDevice:   aws.String(""),
			},
		},
	}

	resp, err := MakeImage(context.TODO(), client, input)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return *resp.ImageId

}
func MakeImage(c context.Context, api EC2CreateImageAPI, input *ec2.CreateImageInput) (*ec2.CreateImageOutput, error) {
	return api.CreateImage(c, input)
}

//https://aws.github.io/aws-sdk-go-v2/docs/
//https://docs.aws.amazon.com/code-samples/latest/catalog/gov2-ec2-CreateInstance-CreateInstancev2.go.html
//https://github.com/gruntwork-io/cloud-nuke/blob/master/aws/ami.go
//https://pkg.go.dev/search?page=5&q=aws-sdk-go-v2

//aws configure
//aws ec2 create-image --instance-id i-07f96abb2dd303e22 --name "My server" --description "An AMI for my server"
//aws ec2 run-instances --image-id ami-0c630a31e852dc15b --count 1 --instance-type t2.micro --key-name keys --security-group-ids sg-0dbcca3589e78cefd
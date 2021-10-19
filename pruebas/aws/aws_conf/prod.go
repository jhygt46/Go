// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

//https://aws.github.io/aws-sdk-go-v2/docs/
//https://docs.aws.amazon.com/code-samples/latest/catalog/gov2-ec2-CreateInstance-CreateInstancev2.go.html

// EC2CreateImageAPI defines the interface for the CreateImage function.
// We use this interface to test the function using a mocked service.
type EC2CreateImageAPI interface {
	CreateImage(ctx context.Context,
		params *ec2.CreateImageInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateImageOutput, error)
}

// MakeImage creates an Amazon Elastic Compute Cloud (Amazon EC2) image.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a CreateImageOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to CreateImage.
func MakeImage(c context.Context, api EC2CreateImageAPI, input *ec2.CreateImageInput) (*ec2.CreateImageOutput, error) {
	return api.CreateImage(c, input)
}

func main() {
	description := flag.String("d", "", "The description of the image")
	instanceID := flag.String("i", "", "The ID of the instance")
	name := flag.String("n", "", "The name of the image")
	flag.Parse()

	if *description == "" || *instanceID == "" || *name == "" {
		fmt.Println("You must supply an image description, instance ID, and image name")
		fmt.Println("(-d IMAGE-DESCRIPTION -i INSTANCE-ID -n IMAGE-NAME")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.CreateImageInput{
		Description: description,
		InstanceId:  instanceID,
		Name:        name,
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
		fmt.Println("Got an error createing image:")
		fmt.Println(err)
		return
	}

	fmt.Println("ID: ", resp.ImageId)
}


func CreateImage(ImageId string) string {
	return "hola"
}

/*
package main

import (
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String("config-bucket-520286683812"),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("first page results:")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}
}
*/

//aws configure
//aws ec2 create-image --instance-id i-07f96abb2dd303e22 --name "My server" --description "An AMI for my server"
//aws ec2 run-instances --image-id ami-0c630a31e852dc15b --count 1 --instance-type t2.micro --key-name keys --security-group-ids sg-0dbcca3589e78cefd
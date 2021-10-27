package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	//"github.com/aws/aws-sdk-go-v2/internal/awsutil"
)

type EC2DeleteImageAPI interface {
	DeregisterImage(ctx context.Context, params *ec2.DeregisterImageInput, optFns ...func(*ec2.Options)) (*ec2.DeregisterImageOutput, error)
}

func main() {

	delete_image("ami-096618d7e0b0294bf")

}

func delete_image(ami string) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.DeregisterImageInput{
		ImageId: aws.String(ami),
	}

	resp, err := DelImage(context.TODO(), client, input)
	if err != nil {
		fmt.Println(err)
	}

	/*
	if awserr := aws.Error(err); awserr != nil {
		fmt.Println("Error:", awserr.Code, awserr.Message)
	} else if err != nil {
		panic(err)
	}

	fmt.Println(awsutil.StringValue(resp))
	*/
	fmt.Println(resp)

}
func DelImage(c context.Context, api EC2DeleteImageAPI, input *ec2.DeregisterImageInput) (*ec2.DeregisterImageOutput, error) {
	return api.DeregisterImage(c, input)
}

//https://pkg.go.dev/github.com/datacratic/aws-sdk-go/service/ec2#DeregisterImageOutput
//https://github.com/gruntwork-io/cloud-nuke/blob/master/aws/ami.go
//https://pkg.go.dev/search?q=aws-sdk-go-v2
//https://aws.github.io/aws-sdk-go-v2/docs/
//https://docs.aws.amazon.com/code-samples/latest/catalog/gov2-ec2-CreateInstance-CreateInstancev2.go.html

//aws configure
//aws ec2 create-image --instance-id i-07f96abb2dd303e22 --name "My server" --description "An AMI for my server"
//aws ec2 run-instances --image-id ami-0c630a31e852dc15b --count 1 --instance-type t2.micro --key-name keys --security-group-ids sg-0dbcca3589e78cefd
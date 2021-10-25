// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2DeleteImageAPI interface {
	TerminateInstances(ctx context.Context, params *ec2.TerminateInstancesInput, optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}

func main() {

	terminate_instance("i-080a8dad14046e40c")

}

func terminate_instance(InstanceId string) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)
	DryRun := true

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{InstanceId},
		DryRun:  &DryRun,
	}

	resp, err := DelInstance(context.TODO(), client, input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(*resp)

}
func DelInstance(c context.Context, api EC2DeleteImageAPI, input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	return api.TerminateInstances(c, input)
}

//https://aws.github.io/aws-sdk-go-v2/docs/
//https://docs.aws.amazon.com/code-samples/latest/catalog/gov2-ec2-CreateInstance-CreateInstancev2.go.html

//aws configure
//aws ec2 create-image --instance-id i-07f96abb2dd303e22 --name "My server" --description "An AMI for my server"
//aws ec2 run-instances --image-id ami-0c630a31e852dc15b --count 1 --instance-type t2.micro --key-name keys --security-group-ids sg-0dbcca3589e78cefd
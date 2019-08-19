package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Bad arguments")
	}
	instanceName := os.Args[1]

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	}))
	svc := ec2.New(sess)

	tagValues := []*string{aws.String(instanceName)}
	filters := []*ec2.Filter{&ec2.Filter{
		Name:   aws.String("tag-value"),
		Values: tagValues,
	}}
	params := &ec2.DescribeInstancesInput{
		Filters:     filters,
	}
	describeInstancesOutput, err := svc.DescribeInstances(params)
	if err != nil {
		log.Fatalf("DescribeInstances Error: %+v\n", err)
	}
	for _, reservation := range describeInstancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" && *tag.Value == instanceName {
					fmt.Println(*instance.PrivateIpAddress)
					os.Exit(0)
				}
			}
		}
	}
	log.Fatal("Instance not found")
}

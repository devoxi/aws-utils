package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/emr"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Bad arguments")
	}
	clusterName := os.Args[1]

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	}))
	svc := emr.New(sess)

	states := []*string{aws.String(emr.ClusterStateWaiting), aws.String(emr.ClusterStateRunning)}
	params := &emr.ListClustersInput{
		ClusterStates: states,
	}
	clustersOutput, err := svc.ListClusters(params)
	if err != nil {
		log.Fatalf("ListCluters Error: %+v\n", err)
	}
	if len(clustersOutput.Clusters) == 0 {
		log.Fatalf("No cluster found")
	}
	for _, element := range clustersOutput.Clusters {
		if *element.Name == clusterName {
			id := element.Id
			groupTypes := []*string{aws.String(emr.InstanceGroupTypeMaster)}
			instancesParams := &emr.ListInstancesInput{
				ClusterId:          id,
				InstanceGroupTypes: groupTypes,
			}
			instancesOutput, err := svc.ListInstances(instancesParams)
			if err != nil {
				log.Fatalf("ListInstances Error: %+v\n", err)
			}
			if len(instancesOutput.Instances) == 0 {
				log.Fatal("No master found")
			}
			master := instancesOutput.Instances[0]
			fmt.Printf("%s\n", *master.PrivateIpAddress)
			os.Exit(0)
		}
	}

	log.Fatal("Cluster not found")
}

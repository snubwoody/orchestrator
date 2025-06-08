package main

import (
	"context"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
)

const projectID = "orchestrator-462314"

func main() {
	zone := "us-east"

	ctx := context.Background()
	fmt.Println("Creating client")

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		fmt.Printf("NewInstancesRESTClient error: %s", err)
	}

	defer instancesClient.Close()

	req := &computepb.ListInstancesRequest{
		Project: projectId,
		Zone:    zone,
	}

	fmt.Println("Fetching instances")
	it := instancesClient.List(ctx, req)
	for {
		instances, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			fmt.Printf("Iter err: %s", err)
			break
		}

		fmt.Printf("%s %s\n", instances.GetName(), instances.GetMachineType())
	}
}

func listZones() {
	ctx := context.Background()
}

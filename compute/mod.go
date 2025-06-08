package compute

import (
	"context"
	"fmt"
	"log"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/snubwoody/orchestrator/config"
	"google.golang.org/api/iterator"
)

func DeleteInstance() {

}

func InsertInstance(zone string) {
	ctx := context.Background()

	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		fmt.Printf("NewInstancesRESTClient error: %s", err)
	}

	defer client.Close()

	instanceName := "my-instance"
	machineType := fmt.Sprintf("zones/%s/machineTypes/e2-micro", zone)
	autoDelete := true
	boot := true
	diskType := "PERSISTENT"
	diskImage := "projects/debian-cloud/global/images/family/debian-11"
	// 10Gb is the minimum size
	var diskSize int64 = 10

	disk := &computepb.AttachedDisk{
		AutoDelete: &autoDelete,
		Boot:       &boot,
		Type:       &diskType,
		InitializeParams: &computepb.AttachedDiskInitializeParams{
			SourceImage: &diskImage,
			DiskSizeGb:  &diskSize,
		},
	}

	networkType := "ONE_TO_ONE_NAT"
	networkName := "External NAT"
	network := "global/networks/default"

	networkInterface := &computepb.NetworkInterface{
		AccessConfigs: []*computepb.AccessConfig{
			{Type: &networkType, Name: &networkName},
		},
		Network: &network,
	}

	resource := &computepb.Instance{
		Name:              &instanceName,
		MachineType:       &machineType,
		Disks:             []*computepb.AttachedDisk{disk},
		NetworkInterfaces: []*computepb.NetworkInterface{networkInterface},
	}

	req := &computepb.InsertInstanceRequest{
		Project:          config.ProjectId,
		Zone:             zone,
		InstanceResource: resource,
	}

	response, err := client.Insert(ctx, req)
	if err != nil {

		log.Fatalf("Error creating instances: %s", err)
		return
	}

	log.Println(response)
}

// List all the active instances in a zone
func ListInstances(zone string) {
	ctx := context.Background()
	fmt.Println("Creating client")

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		fmt.Printf("NewInstancesRESTClient error: %s", err)
	}

	defer instancesClient.Close()

	req := &computepb.ListInstancesRequest{
		Project: config.ProjectId,
		Zone:    zone,
	}

	fmt.Println("Fetching instances")
	it := instancesClient.List(ctx, req)
	count := it.PageInfo().Remaining()

	if count == 0 {
		fmt.Println("No instances found in zone")
		return
	}

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

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

func DeleteInstance(name, zone string) {
	ctx := context.Background()

	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		fmt.Printf("NewInstancesRESTClient error: %s", err)
	}
	defer client.Close()

	req := &computepb.DeleteInstanceRequest{
		Project:  config.ProjectId,
		Zone:     zone,
		Instance: name,
	}

	op, err := client.Delete(ctx, req)
	if err != nil {
		log.Fatalf("Error deleting client: %s", err)
	}

	op.Wait(ctx)
}

func InsertInstance(name, zone string) {
	ctx := context.Background()

	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		fmt.Printf("NewInstancesRESTClient error: %s", err)
	}

	defer client.Close()

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
		Name:              &name,
		MachineType:       &machineType,
		Disks:             []*computepb.AttachedDisk{disk},
		NetworkInterfaces: []*computepb.NetworkInterface{networkInterface},
	}

	req := &computepb.InsertInstanceRequest{
		Project:          config.ProjectId,
		Zone:             zone,
		InstanceResource: resource,
	}

	op, err := client.Insert(ctx, req)
	if err != nil {

		log.Fatalf("Error creating instances: %s", err)
		return
	}

	op.Wait(ctx)
}

// List all the active instances in a zone
func ListInstances(zone string) ([]*computepb.Instance, error) {
	ctx := context.Background()
	fmt.Println("Creating client")

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, err
	}

	defer instancesClient.Close()

	req := &computepb.ListInstancesRequest{
		Project: config.ProjectId,
		Zone:    zone,
	}

	it := instancesClient.List(ctx, req)

	var instances []*computepb.Instance

	for {
		instance, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		instances = append(instances, instance)
	}

	return instances, nil
}

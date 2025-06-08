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

// Compute client
type Client struct {
	context         context.Context
	instancesClient *compute.InstancesClient
}

func NewClient() (Client, error) {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)

	if err != nil {
		return Client{}, err
	}

	client := Client{
		context:         ctx,
		instancesClient: instancesClient,
	}

	return client, err
}

func (c *Client) DeleteInstance(name, zone string) {
	req := &computepb.DeleteInstanceRequest{
		Project:  config.ProjectId,
		Zone:     zone,
		Instance: name,
	}

	op, err := c.instancesClient.Delete(c.context, req)
	if err != nil {
		log.Fatalf("Error deleting client: %s", err)
	}

	op.Wait(c.context)
}

// Deletes an instance without waiting for the operation
// to complete.
func (c *Client) DeleteInstanceAsync(name, zone string) (*compute.Operation, error) {
	req := &computepb.DeleteInstanceRequest{
		Project:  config.ProjectId,
		Zone:     zone,
		Instance: name,
	}

	return c.instancesClient.Delete(c.context, req)
}

func (c *Client) InsertInstance(name, zone string) {
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

	op, err := c.instancesClient.Insert(c.context, req)
	if err != nil {

		log.Fatalf("Error creating instances: %s", err)
		return
	}

	op.Wait(c.context)
}

// List all the active instances in a zone
func (c *Client) ListInstances(zone string) ([]*computepb.Instance, error) {
	fmt.Println("Creating client")

	req := &computepb.ListInstancesRequest{
		Project: config.ProjectId,
		Zone:    zone,
	}

	it := c.instancesClient.List(c.context, req)

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

// Close the client and it's resources
func (c *Client) Close() {
	c.instancesClient.Close()
}
